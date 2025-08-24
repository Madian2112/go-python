# Arquitectura de Microservicios en Python

La arquitectura de microservicios es un enfoque para desarrollar aplicaciones como un conjunto de servicios pequeños e independientes que se ejecutan en su propio proceso y se comunican mediante mecanismos ligeros, a menudo una API HTTP. Python ofrece un ecosistema rico para implementar microservicios gracias a sus frameworks web, herramientas de comunicación asíncrona y bibliotecas para integración de sistemas.

## Principios Fundamentales

### Características de los Microservicios

1. **Desacoplamiento**: Cada servicio es independiente y puede ser desarrollado, desplegado y escalado de forma autónoma.
2. **Especialización**: Cada servicio se enfoca en resolver un problema específico del dominio.
3. **Resiliencia**: El fallo de un servicio no debe afectar a todo el sistema.
4. **Escalabilidad**: Los servicios pueden escalarse de forma independiente según sus necesidades.
5. **Despliegue Independiente**: Cada servicio puede ser desplegado sin afectar a otros servicios.

### Ventajas y Desventajas de Python para Microservicios

**Ventajas:**

1. **Desarrollo Rápido**: Python permite un desarrollo ágil y rápido.
2. **Ecosistema Rico**: Frameworks como Flask, FastAPI y Django REST Framework facilitan la creación de APIs.
3. **Legibilidad**: El código Python es fácil de leer y mantener.
4. **Comunidad Activa**: Gran cantidad de bibliotecas y soporte comunitario.
5. **Soporte para Asincronía**: Con asyncio, Python puede manejar operaciones asíncronas eficientemente.

**Desventajas:**

1. **Rendimiento**: Python puede ser más lento que lenguajes compilados como Go o Rust.
2. **GIL (Global Interpreter Lock)**: Limita la ejecución de múltiples hilos en paralelo.
3. **Consumo de Memoria**: Python tiende a consumir más memoria que otros lenguajes.
4. **Tipado Dinámico**: Puede llevar a errores en tiempo de ejecución en sistemas complejos.

## Estructura Común de un Microservicio en Python

```
/microservice
  /api
    __init__.py
    routes.py
    schemas.py
  /core
    __init__.py
    config.py
    exceptions.py
  /domain
    __init__.py
    models.py
    services.py
  /infrastructure
    __init__.py
    database.py
    repositories.py
    messaging.py
  /tests
    __init__.py
    test_api.py
    test_services.py
  main.py
  Dockerfile
  requirements.txt
```

## Patrones de Diseño para Microservicios

### Arquitectura Hexagonal (Puertos y Adaptadores)

```python
# domain/models.py
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, List

@dataclass
class User:
    id: Optional[str] = None
    username: str
    email: str
    created_at: datetime = datetime.now()
    updated_at: datetime = datetime.now()

# domain/ports.py
from abc import ABC, abstractmethod
from typing import List, Optional
from .models import User

class UserRepository(ABC):
    @abstractmethod
    async def get_by_id(self, user_id: str) -> Optional[User]:
        pass
    
    @abstractmethod
    async def get_all(self) -> List[User]:
        pass
    
    @abstractmethod
    async def create(self, user: User) -> User:
        pass
    
    @abstractmethod
    async def update(self, user: User) -> User:
        pass
    
    @abstractmethod
    async def delete(self, user_id: str) -> bool:
        pass

# domain/services.py
from typing import List, Optional
from .models import User
from .ports import UserRepository

class UserService:
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository
    
    async def get_user(self, user_id: str) -> Optional[User]:
        return await self.user_repository.get_by_id(user_id)
    
    async def get_all_users(self) -> List[User]:
        return await self.user_repository.get_all()
    
    async def create_user(self, user: User) -> User:
        return await self.user_repository.create(user)
    
    async def update_user(self, user: User) -> User:
        return await self.user_repository.update(user)
    
    async def delete_user(self, user_id: str) -> bool:
        return await self.user_repository.delete(user_id)

# infrastructure/repositories.py
from typing import List, Optional
from motor.motor_asyncio import AsyncIOMotorClient
from domain.models import User
from domain.ports import UserRepository

class MongoUserRepository(UserRepository):
    def __init__(self, mongo_client: AsyncIOMotorClient, database_name: str):
        self.db = mongo_client[database_name]
        self.collection = self.db.users
    
    async def get_by_id(self, user_id: str) -> Optional[User]:
        user_data = await self.collection.find_one({"id": user_id})
        if user_data:
            return User(**user_data)
        return None
    
    async def get_all(self) -> List[User]:
        users = []
        async for user_data in self.collection.find():
            users.append(User(**user_data))
        return users
    
    async def create(self, user: User) -> User:
        user_dict = user.__dict__
        await self.collection.insert_one(user_dict)
        return user
    
    async def update(self, user: User) -> User:
        user_dict = user.__dict__
        await self.collection.update_one({"id": user.id}, {"$set": user_dict})
        return user
    
    async def delete(self, user_id: str) -> bool:
        result = await self.collection.delete_one({"id": user_id})
        return result.deleted_count > 0

# api/schemas.py
from pydantic import BaseModel, EmailStr
from datetime import datetime
from typing import Optional

class UserCreate(BaseModel):
    username: str
    email: EmailStr

class UserUpdate(BaseModel):
    username: Optional[str] = None
    email: Optional[EmailStr] = None

class UserResponse(BaseModel):
    id: str
    username: str
    email: EmailStr
    created_at: datetime
    updated_at: datetime

# api/routes.py
from fastapi import APIRouter, Depends, HTTPException, status
from typing import List
from domain.models import User
from domain.services import UserService
from .schemas import UserCreate, UserUpdate, UserResponse
from uuid import uuid4
from datetime import datetime

router = APIRouter()

def get_user_service() -> UserService:
    # En una aplicación real, esto vendría de una fábrica o contenedor de DI
    from infrastructure.repositories import MongoUserRepository
    from motor.motor_asyncio import AsyncIOMotorClient
    
    mongo_client = AsyncIOMotorClient("mongodb://localhost:27017")
    user_repository = MongoUserRepository(mongo_client, "microservice_db")
    return UserService(user_repository)

@router.get("/users/{user_id}", response_model=UserResponse)
async def get_user(user_id: str, user_service: UserService = Depends(get_user_service)):
    user = await user_service.get_user(user_id)
    if not user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    return user

@router.get("/users", response_model=List[UserResponse])
async def get_users(user_service: UserService = Depends(get_user_service)):
    return await user_service.get_all_users()

@router.post("/users", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(user_data: UserCreate, user_service: UserService = Depends(get_user_service)):
    user = User(
        id=str(uuid4()),
        username=user_data.username,
        email=user_data.email
    )
    return await user_service.create_user(user)

@router.put("/users/{user_id}", response_model=UserResponse)
async def update_user(user_id: str, user_data: UserUpdate, user_service: UserService = Depends(get_user_service)):
    existing_user = await user_service.get_user(user_id)
    if not existing_user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    
    # Actualizar solo los campos proporcionados
    if user_data.username is not None:
        existing_user.username = user_data.username
    if user_data.email is not None:
        existing_user.email = user_data.email
    
    existing_user.updated_at = datetime.now()
    
    return await user_service.update_user(existing_user)

@router.delete("/users/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(user_id: str, user_service: UserService = Depends(get_user_service)):
    success = await user_service.delete_user(user_id)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")

# main.py
import uvicorn
from fastapi import FastAPI
from api.routes import router as user_router

app = FastAPI(title="User Microservice")

app.include_router(user_router, prefix="/api/v1")

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
```

### Patrón CQRS (Command Query Responsibility Segregation)

```python
# domain/commands.py
from dataclasses import dataclass
from typing import Optional

@dataclass
class CreateUserCommand:
    username: str
    email: str

@dataclass
class UpdateUserCommand:
    user_id: str
    username: Optional[str] = None
    email: Optional[str] = None

@dataclass
class DeleteUserCommand:
    user_id: str

# domain/queries.py
from dataclasses import dataclass

@dataclass
class GetUserQuery:
    user_id: str

@dataclass
class GetAllUsersQuery:
    limit: int = 100
    offset: int = 0

# domain/handlers.py
from typing import List, Optional
from .models import User
from .commands import CreateUserCommand, UpdateUserCommand, DeleteUserCommand
from .queries import GetUserQuery, GetAllUsersQuery
from .ports import UserRepository
from uuid import uuid4
from datetime import datetime

class CommandHandler:
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository
    
    async def handle_create_user(self, command: CreateUserCommand) -> User:
        user = User(
            id=str(uuid4()),
            username=command.username,
            email=command.email
        )
        return await self.user_repository.create(user)
    
    async def handle_update_user(self, command: UpdateUserCommand) -> Optional[User]:
        user = await self.user_repository.get_by_id(command.user_id)
        if not user:
            return None
        
        if command.username is not None:
            user.username = command.username
        if command.email is not None:
            user.email = command.email
        
        user.updated_at = datetime.now()
        
        return await self.user_repository.update(user)
    
    async def handle_delete_user(self, command: DeleteUserCommand) -> bool:
        return await self.user_repository.delete(command.user_id)

class QueryHandler:
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository
    
    async def handle_get_user(self, query: GetUserQuery) -> Optional[User]:
        return await self.user_repository.get_by_id(query.user_id)
    
    async def handle_get_all_users(self, query: GetAllUsersQuery) -> List[User]:
        # En una implementación real, se pasarían limit y offset al repositorio
        return await self.user_repository.get_all()

# api/routes.py (versión CQRS)
from fastapi import APIRouter, Depends, HTTPException, status
from typing import List
from domain.models import User
from domain.commands import CreateUserCommand, UpdateUserCommand, DeleteUserCommand
from domain.queries import GetUserQuery, GetAllUsersQuery
from domain.handlers import CommandHandler, QueryHandler
from .schemas import UserCreate, UserUpdate, UserResponse

router = APIRouter()

def get_command_handler() -> CommandHandler:
    # En una aplicación real, esto vendría de una fábrica o contenedor de DI
    from infrastructure.repositories import MongoUserRepository
    from motor.motor_asyncio import AsyncIOMotorClient
    
    mongo_client = AsyncIOMotorClient("mongodb://localhost:27017")
    user_repository = MongoUserRepository(mongo_client, "microservice_db")
    return CommandHandler(user_repository)

def get_query_handler() -> QueryHandler:
    # En una aplicación real, esto vendría de una fábrica o contenedor de DI
    from infrastructure.repositories import MongoUserRepository
    from motor.motor_asyncio import AsyncIOMotorClient
    
    mongo_client = AsyncIOMotorClient("mongodb://localhost:27017")
    user_repository = MongoUserRepository(mongo_client, "microservice_db")
    return QueryHandler(user_repository)

@router.get("/users/{user_id}", response_model=UserResponse)
async def get_user(user_id: str, query_handler: QueryHandler = Depends(get_query_handler)):
    query = GetUserQuery(user_id=user_id)
    user = await query_handler.handle_get_user(query)
    if not user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    return user

@router.get("/users", response_model=List[UserResponse])
async def get_users(limit: int = 100, offset: int = 0, query_handler: QueryHandler = Depends(get_query_handler)):
    query = GetAllUsersQuery(limit=limit, offset=offset)
    return await query_handler.handle_get_all_users(query)

@router.post("/users", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(user_data: UserCreate, command_handler: CommandHandler = Depends(get_command_handler)):
    command = CreateUserCommand(username=user_data.username, email=user_data.email)
    return await command_handler.handle_create_user(command)

@router.put("/users/{user_id}", response_model=UserResponse)
async def update_user(user_id: str, user_data: UserUpdate, command_handler: CommandHandler = Depends(get_command_handler)):
    command = UpdateUserCommand(
        user_id=user_id,
        username=user_data.username,
        email=user_data.email
    )
    user = await command_handler.handle_update_user(command)
    if not user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    return user

@router.delete("/users/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(user_id: str, command_handler: CommandHandler = Depends(get_command_handler)):
    command = DeleteUserCommand(user_id=user_id)
    success = await command_handler.handle_delete_user(command)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
```

### API Gateway

```python
# api_gateway/main.py
from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
import httpx
import os
from typing import Dict, Any

app = FastAPI(title="API Gateway")

# Configuración de servicios
SERVICES = {
    "users": os.getenv("USER_SERVICE_URL", "http://user-service:8000"),
    "products": os.getenv("PRODUCT_SERVICE_URL", "http://product-service:8001"),
    "orders": os.getenv("ORDER_SERVICE_URL", "http://order-service:8002"),
}

@app.api_route("/{service}/{path:path}", methods=["GET", "POST", "PUT", "DELETE"])
async def gateway(service: str, path: str, request: Request):
    if service not in SERVICES:
        raise HTTPException(status_code=404, detail=f"Service '{service}' not found")
    
    # Construir la URL del servicio de destino
    target_url = f"{SERVICES[service]}/{path}"
    
    # Obtener el cuerpo de la solicitud
    body = await request.body()
    
    # Obtener los encabezados de la solicitud
    headers = dict(request.headers)
    # Eliminar encabezados específicos de host
    headers.pop("host", None)
    
    # Obtener los parámetros de consulta
    params = dict(request.query_params)
    
    try:
        async with httpx.AsyncClient() as client:
            response = await client.request(
                method=request.method,
                url=target_url,
                headers=headers,
                params=params,
                content=body,
                timeout=30.0
            )
            
            # Devolver la respuesta del servicio
            return JSONResponse(
                content=response.json() if response.content else None,
                status_code=response.status_code,
                headers=dict(response.headers)
            )
    except httpx.RequestError as exc:
        raise HTTPException(status_code=503, detail=f"Error al comunicarse con el servicio: {str(exc)}")

# Endpoint de salud para el API Gateway
@app.get("/health")
async def health():
    return {"status": "ok"}

# Verificar la salud de todos los servicios
@app.get("/health/services")
async def service_health():
    results = {}
    async with httpx.AsyncClient() as client:
        for service_name, service_url in SERVICES.items():
            try:
                response = await client.get(f"{service_url}/health", timeout=5.0)
                results[service_name] = {
                    "status": "up" if response.status_code == 200 else "down",
                    "status_code": response.status_code
                }
            except httpx.RequestError:
                results[service_name] = {"status": "down", "error": "Connection error"}
    
    return results

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8080, reload=True)
```

### Circuit Breaker

```python
# infrastructure/circuit_breaker.py
import asyncio
import time
from enum import Enum
from typing import Callable, Any, TypeVar, Awaitable

T = TypeVar('T')

class CircuitState(Enum):
    CLOSED = 'CLOSED'  # Normal operation, requests pass through
    OPEN = 'OPEN'      # Circuit is open, requests fail fast
    HALF_OPEN = 'HALF_OPEN'  # Testing if service is back online

class CircuitBreaker:
    def __init__(self, failure_threshold: int = 5, recovery_timeout: int = 30, timeout: float = 10.0):
        self.failure_threshold = failure_threshold
        self.recovery_timeout = recovery_timeout
        self.timeout = timeout
        self.state = CircuitState.CLOSED
        self.failure_count = 0
        self.last_failure_time = 0
    
    async def execute(self, func: Callable[..., Awaitable[T]], *args, **kwargs) -> T:
        if self.state == CircuitState.OPEN:
            # Check if recovery timeout has elapsed
            if time.time() - self.last_failure_time > self.recovery_timeout:
                self.state = CircuitState.HALF_OPEN
            else:
                raise CircuitBreakerOpenError("Circuit breaker is open")
        
        try:
            # Set timeout for the function call
            result = await asyncio.wait_for(func(*args, **kwargs), timeout=self.timeout)
            
            # If successful and in HALF_OPEN state, reset the circuit
            if self.state == CircuitState.HALF_OPEN:
                self.reset()
            
            return result
        
        except (asyncio.TimeoutError, Exception) as e:
            # Record the failure
            self.record_failure()
            
            # If we've reached the threshold, open the circuit
            if self.failure_count >= self.failure_threshold:
                self.state = CircuitState.OPEN
                self.last_failure_time = time.time()
            
            # Re-raise the original exception
            if isinstance(e, asyncio.TimeoutError):
                raise CircuitBreakerTimeoutError("Request timed out") from e
            else:
                raise CircuitBreakerError(f"Request failed: {str(e)}") from e
    
    def record_failure(self):
        self.failure_count += 1
    
    def reset(self):
        self.failure_count = 0
        self.state = CircuitState.CLOSED

class CircuitBreakerError(Exception):
    pass

class CircuitBreakerOpenError(CircuitBreakerError):
    pass

class CircuitBreakerTimeoutError(CircuitBreakerError):
    pass

# Ejemplo de uso con un cliente HTTP
class ResilientHttpClient:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.circuit_breaker = CircuitBreaker()
    
    async def get(self, path: str, **kwargs):
        async def _do_request():
            async with httpx.AsyncClient() as client:
                return await client.get(f"{self.base_url}/{path}", **kwargs)
        
        try:
            return await self.circuit_breaker.execute(_do_request)
        except CircuitBreakerOpenError:
            # Implementar fallback o devolver error
            return None
        except CircuitBreakerError as e:
            # Manejar otros errores del circuit breaker
            raise e

# Uso en un servicio
async def get_user_data(user_id: str):
    client = ResilientHttpClient("http://user-service:8000")
    response = await client.get(f"api/v1/users/{user_id}")
    if response:
        return response.json()
    return None
```

### Patrón Saga

```python
# domain/sagas.py
from abc import ABC, abstractmethod
from enum import Enum
from typing import List, Dict, Any, Callable, Awaitable, Optional
import logging

logger = logging.getLogger(__name__)

class SagaStepStatus(Enum):
    NOT_STARTED = "NOT_STARTED"
    STARTED = "STARTED"
    COMPLETED = "COMPLETED"
    FAILED = "FAILED"
    COMPENSATED = "COMPENSATED"

class SagaStep:
    def __init__(self, 
                 name: str,
                 execute_func: Callable[[Dict[str, Any]], Awaitable[Dict[str, Any]]],
                 compensate_func: Callable[[Dict[str, Any]], Awaitable[None]]):
        self.name = name
        self.execute = execute_func
        self.compensate = compensate_func
        self.status = SagaStepStatus.NOT_STARTED
        self.result: Optional[Dict[str, Any]] = None

class Saga:
    def __init__(self, name: str):
        self.name = name
        self.steps: List[SagaStep] = []
        self.context: Dict[str, Any] = {}
    
    def add_step(self, step: SagaStep) -> 'Saga':
        self.steps.append(step)
        return self
    
    async def execute(self) -> Dict[str, Any]:
        current_step_index = 0
        
        try:
            # Execute each step in sequence
            for i, step in enumerate(self.steps):
                current_step_index = i
                logger.info(f"Executing saga step: {step.name}")
                
                step.status = SagaStepStatus.STARTED
                step.result = await step.execute(self.context)
                step.status = SagaStepStatus.COMPLETED
                
                # Update context with step result
                if step.result:
                    self.context.update(step.result)
            
            return self.context
        
        except Exception as e:
            failed_step = self.steps[current_step_index]
            failed_step.status = SagaStepStatus.FAILED
            logger.error(f"Saga step failed: {failed_step.name}. Error: {str(e)}")
            
            # Compensate steps in reverse order
            await self._compensate(current_step_index)
            
            # Re-raise the exception
            raise
    
    async def _compensate(self, failed_step_index: int):
        logger.info(f"Starting compensation for saga: {self.name}")
        
        # Compensate steps in reverse order, up to the failed step
        for i in range(failed_step_index, -1, -1):
            step = self.steps[i]
            
            # Only compensate steps that were completed or the one that failed
            if step.status in [SagaStepStatus.COMPLETED, SagaStepStatus.FAILED]:
                logger.info(f"Compensating saga step: {step.name}")
                
                try:
                    await step.compensate(self.context)
                    step.status = SagaStepStatus.COMPENSATED
                except Exception as e:
                    logger.error(f"Compensation failed for step: {step.name}. Error: {str(e)}")
                    # Continue with other compensations even if one fails

# Ejemplo de uso: Procesamiento de un pedido
async def create_order(user_id: str, product_ids: List[str]):
    # Crear una saga para el proceso de pedido
    order_saga = Saga("CreateOrderSaga")
    
    # Paso 1: Verificar inventario
    async def check_inventory(context):
        # Llamada al servicio de inventario
        inventory_client = ResilientHttpClient("http://inventory-service:8003")
        response = await inventory_client.post("api/v1/inventory/check", json={"product_ids": product_ids})
        
        if not response or response.status_code != 200:
            raise Exception("Failed to check inventory")
        
        inventory_result = response.json()
        if not inventory_result["available"]:
            raise Exception("Products not available in inventory")
        
        return {"inventory_check_id": inventory_result["check_id"]}
    
    async def compensate_inventory(context):
        # No es necesario compensar la verificación de inventario
        pass
    
    # Paso 2: Reservar inventario
    async def reserve_inventory(context):
        inventory_client = ResilientHttpClient("http://inventory-service:8003")
        response = await inventory_client.post("api/v1/inventory/reserve", json={
            "check_id": context["inventory_check_id"],
            "product_ids": product_ids
        })
        
        if not response or response.status_code != 200:
            raise Exception("Failed to reserve inventory")
        
        return {"inventory_reservation_id": response.json()["reservation_id"]}
    
    async def compensate_inventory_reservation(context):
        if "inventory_reservation_id" in context:
            inventory_client = ResilientHttpClient("http://inventory-service:8003")
            await inventory_client.delete(f"api/v1/inventory/reserve/{context['inventory_reservation_id']}")
    
    # Paso 3: Procesar pago
    async def process_payment(context):
        payment_client = ResilientHttpClient("http://payment-service:8004")
        response = await payment_client.post("api/v1/payments", json={
            "user_id": user_id,
            "amount": context.get("total_amount", 100.0),  # En un caso real, se calcularía
            "order_reference": f"ORDER-{context.get('inventory_reservation_id')}"
        })
        
        if not response or response.status_code != 200:
            raise Exception("Payment processing failed")
        
        return {"payment_id": response.json()["payment_id"]}
    
    async def compensate_payment(context):
        if "payment_id" in context:
            payment_client = ResilientHttpClient("http://payment-service:8004")
            await payment_client.post(f"api/v1/payments/{context['payment_id']}/refund")
    
    # Paso 4: Crear pedido
    async def create_order_record(context):
        order_client = ResilientHttpClient("http://order-service:8002")
        response = await order_client.post("api/v1/orders", json={
            "user_id": user_id,
            "product_ids": product_ids,
            "payment_id": context["payment_id"],
            "inventory_reservation_id": context["inventory_reservation_id"]
        })
        
        if not response or response.status_code != 201:
            raise Exception("Failed to create order record")
        
        return {"order_id": response.json()["id"]}
    
    async def compensate_order_creation(context):
        if "order_id" in context:
            order_client = ResilientHttpClient("http://order-service:8002")
            await order_client.delete(f"api/v1/orders/{context['order_id']}")
    
    # Añadir pasos a la saga
    order_saga.add_step(SagaStep("CheckInventory", check_inventory, compensate_inventory))
    order_saga.add_step(SagaStep("ReserveInventory", reserve_inventory, compensate_inventory_reservation))
    order_saga.add_step(SagaStep("ProcessPayment", process_payment, compensate_payment))
    order_saga.add_step(SagaStep("CreateOrderRecord", create_order_record, compensate_order_creation))
    
    # Ejecutar la saga
    try:
        result = await order_saga.execute()
        return {"order_id": result["order_id"], "status": "completed"}
    except Exception as e:
        return {"error": str(e), "status": "failed"}
```

## Comunicación entre Microservicios

### REST

```python
# infrastructure/rest_client.py
import httpx
from typing import Dict, Any, Optional
from .circuit_breaker import CircuitBreaker, CircuitBreakerError

class RestClient:
    def __init__(self, base_url: str, timeout: float = 10.0):
        self.base_url = base_url
        self.timeout = timeout
        self.circuit_breaker = CircuitBreaker()
    
    async def _request(self, method: str, path: str, **kwargs) -> httpx.Response:
        url = f"{self.base_url}/{path}"
        
        async def do_request():
            async with httpx.AsyncClient() as client:
                return await client.request(method, url, timeout=self.timeout, **kwargs)
        
        try:
            return await self.circuit_breaker.execute(do_request)
        except CircuitBreakerError as e:
            raise ServiceCommunicationError(f"Error communicating with service: {str(e)}") from e
    
    async def get(self, path: str, params: Optional[Dict[str, Any]] = None) -> httpx.Response:
        return await self._request("GET", path, params=params)
    
    async def post(self, path: str, json: Optional[Dict[str, Any]] = None) -> httpx.Response:
        return await self._request("POST", path, json=json)
    
    async def put(self, path: str, json: Optional[Dict[str, Any]] = None) -> httpx.Response:
        return await self._request("PUT", path, json=json)
    
    async def delete(self, path: str) -> httpx.Response:
        return await self._request("DELETE", path)

class ServiceCommunicationError(Exception):
    pass

# Ejemplo de uso
class UserServiceClient:
    def __init__(self, base_url: str = "http://user-service:8000"):
        self.client = RestClient(base_url)
    
    async def get_user(self, user_id: str) -> Dict[str, Any]:
        response = await self.client.get(f"api/v1/users/{user_id}")
        if response.status_code == 404:
            return None
        response.raise_for_status()
        return response.json()
    
    async def create_user(self, username: str, email: str) -> Dict[str, Any]:
        response = await self.client.post("api/v1/users", json={"username": username, "email": email})
        response.raise_for_status()
        return response.json()
```

### gRPC

```python
# Definición del servicio en user.proto
"""
syntax = "proto3";

package user;

service UserService {
  rpc GetUser (GetUserRequest) returns (User);
  rpc CreateUser (CreateUserRequest) returns (User);
  rpc UpdateUser (UpdateUserRequest) returns (User);
  rpc DeleteUser (DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers (ListUsersRequest) returns (ListUsersResponse);
}

message GetUserRequest {
  string user_id = 1;
}

message CreateUserRequest {
  string username = 1;
  string email = 2;
}

message UpdateUserRequest {
  string user_id = 1;
  string username = 2;
  string email = 3;
}

message DeleteUserRequest {
  string user_id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

message ListUsersRequest {
  int32 limit = 1;
  int32 offset = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}

message User {
  string id = 1;
  string username = 2;
  string email = 3;
  string created_at = 4;
  string updated_at = 5;
}
"""

# Implementación del servidor gRPC
# infrastructure/grpc_server.py
import grpc
import asyncio
from concurrent import futures
from datetime import datetime
from typing import List, Optional

# Importar los servicios generados por protoc
from proto import user_pb2, user_pb2_grpc

# Importar el servicio de dominio
from domain.services import UserService

class UserServicer(user_pb2_grpc.UserServiceServicer):
    def __init__(self, user_service: UserService):
        self.user_service = user_service
    
    async def GetUser(self, request, context):
        user = await self.user_service.get_user(request.user_id)
        if not user:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"User with ID {request.user_id} not found")
            return user_pb2.User()
        
        return user_pb2.User(
            id=user.id,
            username=user.username,
            email=user.email,
            created_at=user.created_at.isoformat(),
            updated_at=user.updated_at.isoformat()
        )
    
    async def CreateUser(self, request, context):
        from domain.models import User
        from uuid import uuid4
        
        user = User(
            id=str(uuid4()),
            username=request.username,
            email=request.email
        )
        
        created_user = await self.user_service.create_user(user)
        
        return user_pb2.User(
            id=created_user.id,
            username=created_user.username,
            email=created_user.email,
            created_at=created_user.created_at.isoformat(),
            updated_at=created_user.updated_at.isoformat()
        )
    
    async def UpdateUser(self, request, context):
        user = await self.user_service.get_user(request.user_id)
        if not user:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"User with ID {request.user_id} not found")
            return user_pb2.User()
        
        user.username = request.username
        user.email = request.email
        user.updated_at = datetime.now()
        
        updated_user = await self.user_service.update_user(user)
        
        return user_pb2.User(
            id=updated_user.id,
            username=updated_user.username,
            email=updated_user.email,
            created_at=updated_user.created_at.isoformat(),
            updated_at=updated_user.updated_at.isoformat()
        )
    
    async def DeleteUser(self, request, context):
        success = await self.user_service.delete_user(request.user_id)
        if not success:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"User with ID {request.user_id} not found")
        
        return user_pb2.DeleteUserResponse(success=success)
    
    async def ListUsers(self, request, context):
        users = await self.user_service.get_all_users()
        
        response = user_pb2.ListUsersResponse(total=len(users))
        for user in users:
            response.users.append(user_pb2.User(
                id=user.id,
                username=user.username,
                email=user.email,
                created_at=user.created_at.isoformat(),
                updated_at=user.updated_at.isoformat()
            ))
        
        return response

async def serve(user_service: UserService, port: int = 50051):
    server = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServiceServicer_to_server(UserServicer(user_service), server)
    server.add_insecure_port(f'[::]:{port}')
    await server.start()
    print(f"gRPC server started on port {port}")
    await server.wait_for_termination()

# Cliente gRPC
# infrastructure/grpc_client.py
import grpc
from typing import List, Optional, Dict, Any
from proto import user_pb2, user_pb2_grpc

class UserGrpcClient:
    def __init__(self, host: str = "user-service", port: int = 50051):
        self.channel = grpc.aio.insecure_channel(f"{host}:{port}")
        self.stub = user_pb2_grpc.UserServiceStub(self.channel)
    
    async def get_user(self, user_id: str) -> Optional[Dict[str, Any]]:
        try:
            response = await self.stub.GetUser(user_pb2.GetUserRequest(user_id=user_id))
            return {
                "id": response.id,
                "username": response.username,
                "email": response.email,
                "created_at": response.created_at,
                "updated_at": response.updated_at
            }
        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.NOT_FOUND:
                return None
            raise
    
    async def create_user(self, username: str, email: str) -> Dict[str, Any]:
        response = await self.stub.CreateUser(user_pb2.CreateUserRequest(
            username=username,
            email=email
        ))
        
        return {
            "id": response.id,
            "username": response.username,
            "email": response.email,
            "created_at": response.created_at,
            "updated_at": response.updated_at
        }
    
    async def close(self):
        await self.channel.close()

# Uso del cliente gRPC
async def get_user_data(user_id: str):
    client = UserGrpcClient()
    try:
        user = await client.get_user(user_id)
        return user
    finally:
        await client.close()
```

### Mensajería Asíncrona

```python
# infrastructure/messaging.py
import json
import aio_pika
from typing import Dict, Any, Callable, Awaitable
import logging

logger = logging.getLogger(__name__)

class EventPublisher:
    def __init__(self, connection_string: str, exchange_name: str = "microservices"):
        self.connection_string = connection_string
        self.exchange_name = exchange_name
        self.connection = None
        self.channel = None
        self.exchange = None
    
    async def connect(self):
        if self.connection is None or self.connection.is_closed:
            self.connection = await aio_pika.connect_robust(self.connection_string)
            self.channel = await self.connection.channel()
            self.exchange = await self.channel.declare_exchange(
                self.exchange_name,
                aio_pika.ExchangeType.TOPIC,
                durable=True
            )
    
    async def publish(self, routing_key: str, data: Dict[str, Any], headers: Dict[str, str] = None):
        await self.connect()
        
        message_body = json.dumps(data).encode()
        message = aio_pika.Message(
            body=message_body,
            delivery_mode=aio_pika.DeliveryMode.PERSISTENT,
            headers=headers or {}
        )
        
        await self.exchange.publish(message, routing_key=routing_key)
        logger.info(f"Published message to {routing_key}: {data}")
    
    async def close(self):
        if self.connection and not self.connection.is_closed:
            await self.connection.close()
            self.connection = None
            self.channel = None
            self.exchange = None

class EventConsumer:
    def __init__(self, connection_string: str, exchange_name: str = "microservices"):
        self.connection_string = connection_string
        self.exchange_name = exchange_name
        self.connection = None
        self.channel = None
        self.exchange = None
        self.queue = None
    
    async def connect(self, queue_name: str, routing_keys: list):
        if self.connection is None or self.connection.is_closed:
            self.connection = await aio_pika.connect_robust(self.connection_string)
            self.channel = await self.connection.channel()
            
            # Declarar el exchange
            self.exchange = await self.channel.declare_exchange(
                self.exchange_name,
                aio_pika.ExchangeType.TOPIC,
                durable=True
            )
            
            # Declarar la cola
            self.queue = await self.channel.declare_queue(
                queue_name,
                durable=True,
                auto_delete=False
            )
            
            # Vincular la cola a los routing keys
            for routing_key in routing_keys:
                await self.queue.bind(self.exchange, routing_key)
    
    async def consume(self, callback: Callable[[Dict[str, Any], Dict[str, str]], Awaitable[None]]):
        async def process_message(message: aio_pika.IncomingMessage):
            async with message.process():
                try:
                    body = json.loads(message.body.decode())
                    headers = message.headers or {}
                    
                    logger.info(f"Received message: {body}")
                    await callback(body, headers)
                except Exception as e:
                    logger.error(f"Error processing message: {str(e)}")
                    # En producción, podría requerir una estrategia de reintento o cola de errores
        
        await self.queue.consume(process_message)
    
    async def close(self):
        if self.connection and not self.connection.is_closed:
            await self.connection.close()
            self.connection = None
            self.channel = None
            self.exchange = None
            self.queue = None

# Ejemplo de uso: Publicación de eventos
class UserEventPublisher:
    def __init__(self, connection_string: str):
        self.publisher = EventPublisher(connection_string)
    
    async def publish_user_created(self, user_id: str, username: str, email: str):
        await self.publisher.publish(
            routing_key="users.created",
            data={
                "user_id": user_id,
                "username": username,
                "email": email,
                "event_type": "user_created"
            },
            headers={
                "event_id": str(uuid.uuid4()),
                "timestamp": datetime.now().isoformat()
            }
        )
    
    async def publish_user_updated(self, user_id: str, username: str, email: str):
        await self.publisher.publish(
            routing_key="users.updated",
            data={
                "user_id": user_id,
                "username": username,
                "email": email,
                "event_type": "user_updated"
            },
            headers={
                "event_id": str(uuid.uuid4()),
                "timestamp": datetime.now().isoformat()
            }
        )
    
    async def close(self):
        await self.publisher.close()

# Ejemplo de uso: Consumo de eventos
class NotificationService:
    def __init__(self, connection_string: str):
        self.consumer = EventConsumer(connection_string)
    
    async def start(self):
        await self.consumer.connect(
            queue_name="notification_service",
            routing_keys=["users.*", "orders.created"]
        )
        
        await self.consumer.consume(self.process_event)
    
    async def process_event(self, data: Dict[str, Any], headers: Dict[str, str]):
        event_type = data.get("event_type")
        
        if event_type == "user_created":
            await self.send_welcome_email(data["user_id"], data["email"])
        elif event_type == "user_updated":
            await self.send_profile_updated_notification(data["user_id"], data["email"])
        elif event_type == "order_created":
            await self.send_order_confirmation(data["user_id"], data["order_id"])
    
    async def send_welcome_email(self, user_id: str, email: str):
        # Lógica para enviar email de bienvenida
        logger.info(f"Sending welcome email to {email}")
    
    async def send_profile_updated_notification(self, user_id: str, email: str):
        # Lógica para notificar actualización de perfil
        logger.info(f"Sending profile update notification to {email}")
    
    async def send_order_confirmation(self, user_id: str, order_id: str):
        # Lógica para enviar confirmación de pedido
        logger.info(f"Sending order confirmation for order {order_id} to user {user_id}")
    
    async def close(self):
        await self.consumer.close()

# Uso en el servicio de usuarios
async def create_user_with_event(user_data):
    # Crear usuario en la base de datos
    user = await user_service.create_user(user_data)
    
    # Publicar evento
    publisher = UserEventPublisher("amqp://guest:guest@rabbitmq:5672/")
    try:
        await publisher.publish_user_created(
            user_id=user.id,
            username=user.username,
            email=user.email
        )
    finally:
        await publisher.close()
    
    return user
```

## Gestión de Datos

### Base de Datos por Servicio

```python
# infrastructure/database.py
from motor.motor_asyncio import AsyncIOMotorClient
from typing import Dict, Any, List, Optional, TypeVar, Generic, Type
from pydantic import BaseModel
import uuid

T = TypeVar('T', bound=BaseModel)

class MongoRepository(Generic[T]):
    def __init__(self, mongo_client: AsyncIOMotorClient, database_name: str, collection_name: str, model_class: Type[T]):
        self.db = mongo_client[database_name]
        self.collection = self.db[collection_name]
        self.model_class = model_class
    
    async def find_by_id(self, id: str) -> Optional[T]:
        document = await self.collection.find_one({"_id": id})
        if document:
            # Convertir _id a id para compatibilidad con Pydantic
            document["id"] = document.pop("_id")
            return self.model_class(**document)
        return None
    
    async def find_all(self, filter: Dict[str, Any] = None, skip: int = 0, limit: int = 100) -> List[T]:
        cursor = self.collection.find(filter or {})
        cursor = cursor.skip(skip).limit(limit)
        
        result = []
        async for document in cursor:
            # Convertir _id a id para compatibilidad con Pydantic
            document["id"] = document.pop("_id")
            result.append(self.model_class(**document))
        
        return result
    
    async def count(self, filter: Dict[str, Any] = None) -> int:
        return await self.collection.count_documents(filter or {})
    
    async def create(self, model: T) -> T:
        # Asegurar que el modelo tiene un ID
        data = model.dict()
        if "id" not in data or not data["id"]:
            data["id"] = str(uuid.uuid4())
        
        # Convertir id a _id para MongoDB
        data["_id"] = data.pop("id")
        
        await self.collection.insert_one(data)
        
        # Restaurar id para devolver el modelo
        data["id"] = data.pop("_id")
        return self.model_class(**data)
    
    async def update(self, id: str, model: T) -> Optional[T]:
        data = model.dict(exclude={"id"})  # Excluir id para no sobrescribirlo
        
        result = await self.collection.update_one(
            {"_id": id},
            {"$set": data}
        )
        
        if result.matched_count > 0:
            return await self.find_by_id(id)
        return None
    
    async def delete(self, id: str) -> bool:
        result = await self.collection.delete_one({"_id": id})
        return result.deleted_count > 0

# Ejemplo de uso con SQLAlchemy para otro servicio
# infrastructure/sql_database.py
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker, declarative_base
from sqlalchemy import Column, String, DateTime, Integer, ForeignKey
from sqlalchemy.future import select
from datetime import datetime
import uuid
from typing import List, Optional, TypeVar, Generic, Type, Any, Dict

Base = declarative_base()

class Product(Base):
    __tablename__ = "products"
    
    id = Column(String, primary_key=True)
    name = Column(String, nullable=False)
    description = Column(String)
    price = Column(Integer, nullable=False)  # Precio en centavos
    stock = Column(Integer, nullable=False, default=0)
    created_at = Column(DateTime, default=datetime.now)
    updated_at = Column(DateTime, default=datetime.now, onupdate=datetime.now)

class SQLRepository(Generic[T]):
    def __init__(self, engine, model_class: Type[Base]):
        self.engine = engine
        self.model_class = model_class
        self.session_factory = sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)
    
    async def find_by_id(self, id: str) -> Optional[T]:
        async with self.session_factory() as session:
            result = await session.execute(select(self.model_class).filter_by(id=id))
            return result.scalars().first()
    
    async def find_all(self, skip: int = 0, limit: int = 100, **filters) -> List[T]:
        async with self.session_factory() as session:
            query = select(self.model_class).filter_by(**filters).offset(skip).limit(limit)
            result = await session.execute(query)
            return result.scalars().all()
    
    async def create(self, **data) -> T:
        if "id" not in data or not data["id"]:
            data["id"] = str(uuid.uuid4())
        
        instance = self.model_class(**data)
        
        async with self.session_factory() as session:
            session.add(instance)
            await session.commit()
            await session.refresh(instance)
            return instance
    
    async def update(self, id: str, **data) -> Optional[T]:
        async with self.session_factory() as session:
            instance = await self.find_by_id(id)
            if not instance:
                return None
            
            for key, value in data.items():
                setattr(instance, key, value)
            
            session.add(instance)
            await session.commit()
            await session.refresh(instance)
            return instance
    
    async def delete(self, id: str) -> bool:
        async with self.session_factory() as session:
            instance = await self.find_by_id(id)
            if not instance:
                return False
            
            await session.delete(instance)
            await session.commit()
            return True

# Configuración de la base de datos
async def init_db():
    # SQLite para desarrollo, en producción usaríamos PostgreSQL
    engine = create_async_engine("sqlite+aiosqlite:///./product_service.db")
    
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    
    return engine
```

### Consultas entre Servicios

```python
# domain/services.py
from typing import List, Optional, Dict, Any
from .models import Order, OrderItem
from .ports import OrderRepository
from infrastructure.rest_client import RestClient, ServiceCommunicationError

class OrderService:
    def __init__(self, order_repository: OrderRepository, product_service_url: str, user_service_url: str):
        self.order_repository = order_repository
        self.product_client = RestClient(product_service_url)
        self.user_client = RestClient(user_service_url)
    
    async def create_order(self, user_id: str, items: List[Dict[str, Any]]) -> Optional[Order]:
        # Verificar que el usuario existe
        try:
            user_response = await self.user_client.get(f"api/v1/users/{user_id}")
            if user_response.status_code != 200:
                raise ValueError(f"Usuario con ID {user_id} no encontrado")
            
            # Verificar que los productos existen y tienen stock suficiente
            product_ids = [item["product_id"] for item in items]
            products_response = await self.product_client.post("api/v1/products/check", json={
                "product_ids": product_ids,
                "quantities": {item["product_id"]: item["quantity"] for item in items}
            })
            
            if products_response.status_code != 200:
                raise ValueError("Uno o más productos no están disponibles")
            
            # Calcular el total del pedido
            products_data = products_response.json()["products"]
            total_amount = sum(products_data[item["product_id"]]["price"] * item["quantity"] for item in items)
            
            # Crear el pedido
            order_items = [OrderItem(
                product_id=item["product_id"],
                quantity=item["quantity"],
                price=products_data[item["product_id"]]["price"]
            ) for item in items]
            
            order = Order(
                user_id=user_id,
                items=order_items,
                total_amount=total_amount,
                status="pending"
            )
            
            # Guardar el pedido
            created_order = await self.order_repository.create(order)
            
            # Actualizar el inventario (esto podría hacerse de forma asíncrona con eventos)
            await self.product_client.post("api/v1/products/reserve", json={
                "order_id": created_order.id,
                "items": [{
                    "product_id": item.product_id,
                    "quantity": item.quantity
                } for item in created_order.items]
            })
            
            return created_order
        
        except ServiceCommunicationError as e:
            # Manejar errores de comunicación
            raise ValueError(f"Error al comunicarse con otros servicios: {str(e)}")

### Consistencia Eventual

```python
# infrastructure/event_sourcing.py
from typing import Dict, Any, List, Optional, Type, Generic, TypeVar
from pydantic import BaseModel
from datetime import datetime
import uuid
import json

T = TypeVar('T', bound=BaseModel)

class Event(BaseModel):
    id: str
    aggregate_id: str
    aggregate_type: str
    event_type: str
    data: Dict[str, Any]
    metadata: Dict[str, Any]
    created_at: datetime

class EventStore:
    def __init__(self, mongo_client, database_name: str, collection_name: str = "events"):
        self.db = mongo_client[database_name]
        self.collection = self.db[collection_name]
    
    async def save_event(self, event: Event) -> Event:
        event_dict = event.dict()
        await self.collection.insert_one(event_dict)
        return event
    
    async def get_events_by_aggregate(self, aggregate_id: str, aggregate_type: str) -> List[Event]:
        cursor = self.collection.find({
            "aggregate_id": aggregate_id,
            "aggregate_type": aggregate_type
        }).sort("created_at", 1)  # Ordenar por fecha de creación
        
        events = []
        async for doc in cursor:
            events.append(Event(**doc))
        
        return events

class Aggregate(Generic[T]):
    def __init__(self, model_class: Type[T], event_store: EventStore):
        self.model_class = model_class
        self.event_store = event_store
    
    async def load(self, aggregate_id: str) -> Optional[T]:
        events = await self.event_store.get_events_by_aggregate(
            aggregate_id=aggregate_id,
            aggregate_type=self.model_class.__name__
        )
        
        if not events:
            return None
        
        # Reconstruir el estado aplicando eventos
        return self._apply_events(events)
    
    def _apply_events(self, events: List[Event]) -> T:
        # Crear una instancia vacía
        instance = self.model_class(id=events[0].aggregate_id)
        
        # Aplicar cada evento
        for event in events:
            self._apply_event(instance, event)
        
        return instance
    
    def _apply_event(self, instance: T, event: Event) -> None:
        # Implementación específica para cada tipo de evento
        method_name = f"apply_{event.event_type}"
        if hasattr(self, method_name):
            method = getattr(self, method_name)
            method(instance, event.data)

# Ejemplo de uso con un agregado de Pedido
class OrderAggregate(Aggregate[Order]):
    def apply_order_created(self, order: Order, data: Dict[str, Any]) -> None:
        order.user_id = data["user_id"]
        order.total_amount = data["total_amount"]
        order.status = data["status"]
        order.created_at = datetime.fromisoformat(data["created_at"])
        order.updated_at = datetime.fromisoformat(data["updated_at"])
        
        # Reconstruir items
        order.items = [OrderItem(**item) for item in data["items"]]
    
    def apply_order_paid(self, order: Order, data: Dict[str, Any]) -> None:
        order.status = "paid"
        order.payment_id = data["payment_id"]
        order.updated_at = datetime.fromisoformat(data["updated_at"])
    
    def apply_order_shipped(self, order: Order, data: Dict[str, Any]) -> None:
        order.status = "shipped"
        order.tracking_number = data["tracking_number"]
        order.shipped_at = datetime.fromisoformat(data["shipped_at"])
        order.updated_at = datetime.fromisoformat(data["updated_at"])
    
    def apply_order_delivered(self, order: Order, data: Dict[str, Any]) -> None:
        order.status = "delivered"
        order.delivered_at = datetime.fromisoformat(data["delivered_at"])
        order.updated_at = datetime.fromisoformat(data["updated_at"])
    
    def apply_order_cancelled(self, order: Order, data: Dict[str, Any]) -> None:
        order.status = "cancelled"
        order.cancellation_reason = data.get("reason")
        order.updated_at = datetime.fromisoformat(data["updated_at"])

# Servicio que utiliza Event Sourcing
class EventSourcedOrderService:
    def __init__(self, event_store: EventStore, publisher: EventPublisher):
        self.event_store = event_store
        self.order_aggregate = OrderAggregate(Order, event_store)
        self.publisher = publisher
    
    async def create_order(self, user_id: str, items: List[Dict[str, Any]], total_amount: int) -> Order:
        # Crear un nuevo ID para el pedido
        order_id = str(uuid.uuid4())
        
        # Crear el evento
        now = datetime.now()
        event = Event(
            id=str(uuid.uuid4()),
            aggregate_id=order_id,
            aggregate_type="Order",
            event_type="order_created",
            data={
                "user_id": user_id,
                "items": items,
                "total_amount": total_amount,
                "status": "pending",
                "created_at": now.isoformat(),
                "updated_at": now.isoformat()
            },
            metadata={
                "user_id": user_id,
                "source": "order_service"
            },
            created_at=now
        )
        
        # Guardar el evento
        await self.event_store.save_event(event)
        
        # Publicar el evento
        await self.publisher.publish(
            routing_key="orders.created",
            data={
                "order_id": order_id,
                "user_id": user_id,
                "total_amount": total_amount,
                "items": items,
                "event_type": "order_created"
            }
        )
        
        # Cargar y devolver el pedido
        return await self.order_aggregate.load(order_id)
    
    async def mark_order_as_paid(self, order_id: str, payment_id: str) -> Order:
        # Verificar que el pedido existe
        order = await self.order_aggregate.load(order_id)
        if not order:
            raise ValueError(f"Pedido con ID {order_id} no encontrado")
        
        if order.status != "pending":
            raise ValueError(f"El pedido no está en estado pendiente, estado actual: {order.status}")
        
        # Crear el evento
        now = datetime.now()
        event = Event(
            id=str(uuid.uuid4()),
            aggregate_id=order_id,
            aggregate_type="Order",
            event_type="order_paid",
            data={
                "payment_id": payment_id,
                "updated_at": now.isoformat()
            },
            metadata={
                "source": "payment_service"
            },
            created_at=now
        )
        
        # Guardar el evento
        await self.event_store.save_event(event)
        
        # Publicar el evento
        await self.publisher.publish(
            routing_key="orders.paid",
            data={
                "order_id": order_id,
                "payment_id": payment_id,
                "event_type": "order_paid"
            }
        )
        
        # Cargar y devolver el pedido actualizado
        return await self.order_aggregate.load(order_id)
```

## Seguridad

### Autenticación y Autorización

```python
# infrastructure/security.py
import jwt
from datetime import datetime, timedelta
from typing import Dict, Any, Optional, List
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
import bcrypt
import os

# Configuración
SECRET_KEY = os.getenv("JWT_SECRET_KEY", "your-secret-key")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

class JWTService:
    @staticmethod
    def create_access_token(data: Dict[str, Any], expires_delta: Optional[timedelta] = None) -> str:
        to_encode = data.copy()
        
        if expires_delta:
            expire = datetime.utcnow() + expires_delta
        else:
            expire = datetime.utcnow() + timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
        
        to_encode.update({"exp": expire})
        encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
        return encoded_jwt
    
    @staticmethod
    def decode_token(token: str) -> Dict[str, Any]:
        try:
            payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
            return payload
        except jwt.PyJWTError:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Could not validate credentials",
                headers={"WWW-Authenticate": "Bearer"},
            )

class PasswordService:
    @staticmethod
    def hash_password(password: str) -> str:
        salt = bcrypt.gensalt()
        hashed = bcrypt.hashpw(password.encode(), salt)
        return hashed.decode()
    
    @staticmethod
    def verify_password(plain_password: str, hashed_password: str) -> bool:
        return bcrypt.checkpw(plain_password.encode(), hashed_password.encode())

# Middleware de autenticación
async def get_current_user(token: str = Depends(oauth2_scheme)) -> Dict[str, Any]:
    payload = JWTService.decode_token(token)
    user_id = payload.get("sub")
    if user_id is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )
    return {"id": user_id, "roles": payload.get("roles", [])}

# Middleware de autorización
def has_role(required_roles: List[str]):
    async def role_checker(current_user: Dict[str, Any] = Depends(get_current_user)):
        user_roles = current_user.get("roles", [])
        for role in required_roles:
            if role in user_roles:
                return current_user
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions"
        )
    return role_checker

# Ejemplo de uso en rutas
@router.get("/users/me")
async def read_users_me(current_user: Dict[str, Any] = Depends(get_current_user)):
    return current_user

@router.get("/admin/users")
async def get_all_users(current_user: Dict[str, Any] = Depends(has_role(["admin"]))):
    # Solo los administradores pueden acceder a esta ruta
    return {"message": "Lista de todos los usuarios"}
```

### Seguridad TLS

```python
# main.py con configuración TLS
import uvicorn
from fastapi import FastAPI
from api.routes import router as user_router
import ssl

app = FastAPI(title="User Microservice")

app.include_router(user_router, prefix="/api/v1")

if __name__ == "__main__":
    # Configurar contexto SSL
    ssl_context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    ssl_context.load_cert_chain("./certs/server.crt", keyfile="./certs/server.key")
    
    uvicorn.run(
        "main:app", 
        host="0.0.0.0", 
        port=8443, 
        ssl_keyfile="./certs/server.key",
        ssl_certfile="./certs/server.crt"
    )
```

### Gestión de Secretos

```python
# infrastructure/secrets.py
import hvac
from typing import Dict, Any, Optional
import os

class VaultClient:
    def __init__(self, url: str = None, token: str = None):
        self.url = url or os.getenv("VAULT_ADDR", "http://vault:8200")
        self.token = token or os.getenv("VAULT_TOKEN")
        self.client = hvac.Client(url=self.url, token=self.token)
    
    def is_authenticated(self) -> bool:
        return self.client.is_authenticated()
    
    def get_secret(self, path: str, key: Optional[str] = None) -> Any:
        try:
            response = self.client.secrets.kv.v2.read_secret_version(path=path)
            data = response["data"]["data"]
            
            if key:
                return data.get(key)
            return data
        except Exception as e:
            raise ValueError(f"Error retrieving secret: {str(e)}")
    
    def set_secret(self, path: str, data: Dict[str, Any]) -> None:
        try:
            self.client.secrets.kv.v2.create_or_update_secret(
                path=path,
                secret=data
            )
        except Exception as e:
            raise ValueError(f"Error setting secret: {str(e)}")

# Ejemplo de uso
class DatabaseConfig:
    def __init__(self, vault_client: VaultClient):
        self.vault_client = vault_client
    
    def get_connection_string(self, service_name: str) -> str:
        db_secrets = self.vault_client.get_secret(f"database/{service_name}")
        
        host = db_secrets["host"]
        port = db_secrets["port"]
        username = db_secrets["username"]
        password = db_secrets["password"]
        database = db_secrets["database"]
        
        return f"mongodb://{username}:{password}@{host}:{port}/{database}"

# Uso en la aplicación
from fastapi import FastAPI, Depends

app = FastAPI()

def get_vault_client():
    client = VaultClient()
    if not client.is_authenticated():
        raise Exception("Vault client not authenticated")
    return client

def get_db_config(vault_client: VaultClient = Depends(get_vault_client)):
    return DatabaseConfig(vault_client)

@app.on_event("startup")
async def startup_db_client():
    app.vault_client = VaultClient()
    app.db_config = DatabaseConfig(app.vault_client)
    
    # Obtener cadena de conexión
    connection_string = app.db_config.get_connection_string("user-service")
    
    # Configurar cliente de base de datos
    from motor.motor_asyncio import AsyncIOMotorClient
    app.mongodb_client = AsyncIOMotorClient(connection_string)
    app.mongodb = app.mongodb_client["user_service"]

@app.on_event("shutdown")
async def shutdown_db_client():
    app.mongodb_client.close()
```

## Observabilidad

### Logging

```python
# infrastructure/logging.py
import logging
import json
from datetime import datetime
from typing import Dict, Any, Optional
import sys
import uuid
from contextvars import ContextVar

# Variable de contexto para el ID de correlación
correlation_id: ContextVar[str] = ContextVar('correlation_id', default='')

class JsonFormatter(logging.Formatter):
    def format(self, record):
        log_record = {
            "timestamp": datetime.utcnow().isoformat(),
            "level": record.levelname,
            "message": record.getMessage(),
            "module": record.module,
            "function": record.funcName,
            "line": record.lineno,
        }
        
        # Añadir ID de correlación si existe
        corr_id = correlation_id.get()
        if corr_id:
            log_record["correlation_id"] = corr_id
        
        # Añadir excepción si existe
        if record.exc_info:
            log_record["exception"] = self.formatException(record.exc_info)
        
        # Añadir datos adicionales si existen
        if hasattr(record, 'data'):
            log_record.update(record.data)
        
        return json.dumps(log_record)

def setup_logging(service_name: str, log_level: str = "INFO"):
    logger = logging.getLogger(service_name)
    logger.setLevel(getattr(logging, log_level))
    
    # Crear handler para stdout
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(JsonFormatter())
    logger.addHandler(handler)
    
    return logger

class Logger:
    def __init__(self, service_name: str, log_level: str = "INFO"):
        self.logger = setup_logging(service_name, log_level)
    
    def set_correlation_id(self, corr_id: Optional[str] = None):
        correlation_id.set(corr_id or str(uuid.uuid4()))
    
    def get_correlation_id(self) -> str:
        return correlation_id.get()
    
    def _log(self, level: str, message: str, data: Optional[Dict[str, Any]] = None):
        extra = {"data": data or {}}
        getattr(self.logger, level.lower())(message, extra=extra)
    
    def debug(self, message: str, data: Optional[Dict[str, Any]] = None):
        self._log("DEBUG", message, data)
    
    def info(self, message: str, data: Optional[Dict[str, Any]] = None):
        self._log("INFO", message, data)
    
    def warning(self, message: str, data: Optional[Dict[str, Any]] = None):
        self._log("WARNING", message, data)
    
    def error(self, message: str, data: Optional[Dict[str, Any]] = None):
        self._log("ERROR", message, data)
    
    def critical(self, message: str, data: Optional[Dict[str, Any]] = None):
        self._log("CRITICAL", message, data)

# Middleware para FastAPI que añade ID de correlación
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware

class CorrelationIdMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, logger: Logger):
        super().__init__(app)
        self.logger = logger
    
    async def dispatch(self, request: Request, call_next):
        # Obtener ID de correlación del encabezado o generar uno nuevo
        corr_id = request.headers.get("X-Correlation-ID") or str(uuid.uuid4())
        self.logger.set_correlation_id(corr_id)
        
        # Registrar la solicitud entrante
        self.logger.info(
            f"Incoming request: {request.method} {request.url.path}",
            {"http": {"method": request.method, "url": str(request.url)}}
        )
        
        # Procesar la solicitud
        response = await call_next(request)
        
        # Añadir ID de correlación a la respuesta
        response.headers["X-Correlation-ID"] = corr_id
        
        # Registrar la respuesta saliente
        self.logger.info(
            f"Outgoing response: {response.status_code}",
            {"http": {"status_code": response.status_code}}
        )
        
        return response

# Uso en la aplicación
from fastapi import FastAPI

app = FastAPI()

# Configurar logger
logger = Logger("user-service")

# Añadir middleware
app.add_middleware(CorrelationIdMiddleware, logger=logger)

# Ejemplo de uso en un endpoint
@app.get("/api/v1/users/{user_id}")
async def get_user(user_id: str):
    logger.info(f"Fetching user", {"user_id": user_id})
    
    # Lógica para obtener usuario
    
    logger.info(f"User fetched successfully", {"user_id": user_id})
    return {"id": user_id, "username": "example"}
```

### Métricas

```python
# infrastructure/metrics.py
from prometheus_client import Counter, Histogram, Gauge, Summary, CollectorRegistry, push_to_gateway, generate_latest
from typing import Dict, Any, Optional, Callable
import time
from fastapi import FastAPI, Request, Response
from starlette.middleware.base import BaseHTTPMiddleware

class PrometheusMetrics:
    def __init__(self, app: FastAPI, service_name: str, push_gateway_url: Optional[str] = None):
        self.app = app
        self.service_name = service_name
        self.push_gateway_url = push_gateway_url
        self.registry = CollectorRegistry()
        
        # Métricas HTTP
        self.http_requests_total = Counter(
            'http_requests_total',
            'Total number of HTTP requests',
            ['method', 'endpoint', 'status'],
            registry=self.registry
        )
        
        self.http_request_duration_seconds = Histogram(
            'http_request_duration_seconds',
            'HTTP request duration in seconds',
            ['method', 'endpoint'],
            registry=self.registry
        )
        
        self.http_requests_in_progress = Gauge(
            'http_requests_in_progress',
            'Number of HTTP requests in progress',
            ['method', 'endpoint'],
            registry=self.registry
        )
        
        # Métricas de dependencias
        self.dependency_request_duration_seconds = Histogram(
            'dependency_request_duration_seconds',
            'Dependency request duration in seconds',
            ['dependency_name', 'operation', 'status'],
            registry=self.registry
        )
        
        self.dependency_errors_total = Counter(
            'dependency_errors_total',
            'Total number of dependency errors',
            ['dependency_name', 'operation', 'error_type'],
            registry=self.registry
        )
        
        # Métricas de negocio
        self.business_operation_duration_seconds = Histogram(
            'business_operation_duration_seconds',
            'Business operation duration in seconds',
            ['operation'],
            registry=self.registry
        )
        
        # Añadir middleware para métricas HTTP
        app.add_middleware(PrometheusMiddleware, metrics=self)
        
        # Añadir endpoint para exponer métricas
        @app.get('/metrics')
        async def metrics():
            if self.push_gateway_url:
                push_to_gateway(self.push_gateway_url, job=self.service_name, registry=self.registry)
            return Response(content=generate_latest(self.registry), media_type="text/plain")
    
    def track_dependency(self, dependency_name: str, operation: str):
        def decorator(func):
            async def wrapper(*args, **kwargs):
                start_time = time.time()
                status = "success"
                try:
                    result = await func(*args, **kwargs)
                    return result
                except Exception as e:
                    status = "error"
                    self.dependency_errors_total.labels(
                        dependency_name=dependency_name,
                        operation=operation,
                        error_type=type(e).__name__
                    ).inc()
                    raise
                finally:
                    self.dependency_request_duration_seconds.labels(
                        dependency_name=dependency_name,
                        operation=operation,
                        status=status
                    ).observe(time.time() - start_time)
            return wrapper
        return decorator
    
    def track_business_operation(self, operation: str):
        def decorator(func):
            async def wrapper(*args, **kwargs):
                start_time = time.time()
                try:
                    result = await func(*args, **kwargs)
                    return result
                finally:
                    self.business_operation_duration_seconds.labels(
                        operation=operation
                    ).observe(time.time() - start_time)
            return wrapper
        return decorator

class PrometheusMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, metrics: PrometheusMetrics):
        super().__init__(app)
        self.metrics = metrics
    
    async def dispatch(self, request: Request, call_next):
        method = request.method
        path = request.url.path
        
        if path == "/metrics":
            return await call_next(request)
        
        # Incrementar contador de solicitudes en progreso
        self.metrics.http_requests_in_progress.labels(
            method=method,
            endpoint=path
        ).inc()
        
        start_time = time.time()
        
        try:
            response = await call_next(request)
            status_code = response.status_code
            return response
        except Exception as e:
            status_code = 500
            raise
        finally:
            # Decrementar contador de solicitudes en progreso
            self.metrics.http_requests_in_progress.labels(
                method=method,
                endpoint=path
            ).dec()
            
            # Registrar duración
            self.metrics.http_request_duration_seconds.labels(
                method=method,
                endpoint=path
            ).observe(time.time() - start_time)
            
            # Incrementar contador total
            self.metrics.http_requests_total.labels(
                method=method,
                endpoint=path,
                status=status_code
            ).inc()

# Ejemplo de uso
from fastapi import FastAPI

app = FastAPI()
metrics = PrometheusMetrics(app, "user-service")

# Ejemplo de uso del decorador para dependencias
@metrics.track_dependency("mongodb", "get_user")
async def get_user_from_db(user_id: str):
    # Lógica para obtener usuario de la base de datos
    return {"id": user_id, "username": "example"}

# Ejemplo de uso del decorador para operaciones de negocio
@app.get("/api/v1/users/{user_id}")
@metrics.track_business_operation("get_user")
async def get_user(user_id: str):
    return await get_user_from_db(user_id)
```

### Trazabilidad

```python
# infrastructure/tracing.py
from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.aiohttp import AioHttpClientInstrumentor
from opentelemetry.instrumentation.asyncpg import AsyncPGInstrumentor
from opentelemetry.instrumentation.motor import MotorInstrumentor
from fastapi import FastAPI, Request
from typing import Dict, Any, Optional, Callable
import functools

def setup_tracing(app: FastAPI, service_name: str, jaeger_host: str = "jaeger", jaeger_port: int = 6831):
    # Configurar el proveedor de trazas
    resource = Resource(attributes={SERVICE_NAME: service_name})
    provider = TracerProvider(resource=resource)
    
    # Configurar el exportador Jaeger
    jaeger_exporter = JaegerExporter(
        agent_host_name=jaeger_host,
        agent_port=jaeger_port,
    )
    
    # Añadir el procesador de spans
    processor = BatchSpanProcessor(jaeger_exporter)
    provider.add_span_processor(processor)
    
    # Establecer el proveedor global
    trace.set_tracer_provider(provider)
    
    # Obtener un tracer para el servicio
    tracer = trace.get_tracer(service_name)
    
    # Instrumentar FastAPI
    FastAPIInstrumentor.instrument_app(app)
    
    # Instrumentar clientes HTTP
    AioHttpClientInstrumentor().instrument()
    
    # Instrumentar bases de datos
    AsyncPGInstrumentor().instrument()
    MotorInstrumentor().instrument()
    
    return tracer

def trace_function(name: Optional[str] = None, attributes: Optional[Dict[str, Any]] = None):
    def decorator(func):
        @functools.wraps(func)
        async def wrapper(*args, **kwargs):
            # Obtener el tracer
            tracer = trace.get_tracer(__name__)
            
            # Nombre del span
            span_name = name or func.__name__
            
            # Crear un nuevo span
            with tracer.start_as_current_span(span_name) as span:
                # Añadir atributos al span
                if attributes:
                    for key, value in attributes.items():
                        span.set_attribute(key, value)
                
                # Ejecutar la función
                try:
                    result = await func(*args, **kwargs)
                    return result
                except Exception as e:
                    # Registrar la excepción en el span
                    span.record_exception(e)
                    span.set_status(trace.Status(trace.StatusCode.ERROR, str(e)))
                    raise
        
        return wrapper
    return decorator

# Ejemplo de uso
from fastapi import FastAPI, Depends

app = FastAPI()
tracer = setup_tracing(app, "user-service")

# Ejemplo de uso del decorador para funciones
@trace_function(name="get_user_from_db", attributes={"db.type": "mongodb"})
async def get_user_from_db(user_id: str):
    # Lógica para obtener usuario de la base de datos
    return {"id": user_id, "username": "example"}

# Ejemplo de uso en un endpoint
@app.get("/api/v1/users/{user_id}")
async def get_user(user_id: str):
    # El endpoint ya está instrumentado por FastAPIInstrumentor
    # Podemos añadir spans personalizados dentro del endpoint
    
    current_span = trace.get_current_span()
    current_span.set_attribute("user.id", user_id)
    
    # Llamar a la función instrumentada
    user = await get_user_from_db(user_id)
    
    return user
```

## Despliegue de Microservicios

### Contenedores con Docker

```dockerfile
# Dockerfile
FROM python:3.9-slim

WORKDIR /app

# Instalar dependencias
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copiar código fuente
COPY . .

# Exponer puerto
EXPOSE 8000

# Comando para ejecutar la aplicación
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

### Orquestación con Kubernetes

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: yourusername/user-service:latest
        ports:
        - containerPort: 8000
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: mongodb_uri
        - name: JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: jwt_secret
        - name: RABBITMQ_URI
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: rabbitmq_uri
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 80
    targetPort: 8000
  type: ClusterIP
---
apiVersion: v1
kind: Secret
metadata:
  name: user-service-secrets
type: Opaque
data:
  mongodb_uri: bW9uZ29kYjovL3VzZXJuYW1lOnBhc3N3b3JkQG1vbmdvZGI6MjcwMTcvdXNlcl9zZXJ2aWNl
  jwt_secret: c3VwZXJfc2VjcmV0X2tleQ==
  rabbitmq_uri: YW1xcDovL2d1ZXN0Omd1ZXN0QHJhYmJpdG1xOjU2NzIv
```

### Despliegue Continuo

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: 3.9

    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install pytest pytest-asyncio
        if [ -f requirements.txt ]; then pip install -r requirements.txt; fi

    - name: Test with pytest
      run: |
        pytest

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v2

    - name: Build and push Docker image
      uses: docker/build-push-action@v2
      with:
        context: .
        push: true
        tags: yourusername/user-service:latest

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Set up kubectl
      uses: azure/setup-kubectl@v1

    - name: Set Kubernetes context
      uses: azure/k8s-set-context@v1
      with:
        kubeconfig: ${{ secrets.KUBE_CONFIG }}

    - name: Deploy to Kubernetes
      run: |
        kubectl apply -f kubernetes/deployment.yaml
        kubectl rollout restart deployment/user-service
```

## Conclusiones

La arquitectura de microservicios en Python ofrece numerosas ventajas para sistemas distribuidos, incluyendo:

1. **Escalabilidad**: Cada servicio puede escalarse de forma independiente según sus necesidades.
2. **Resiliencia**: El fallo de un servicio no afecta a todo el sistema.
3. **Tecnología heterogénea**: Cada servicio puede utilizar la tecnología más adecuada para su función.
4. **Despliegue independiente**: Los servicios pueden desplegarse sin afectar al resto del sistema.
5. **Equipos autónomos**: Diferentes equipos pueden trabajar en diferentes servicios.

Sin embargo, también presenta desafíos:

1. **Complejidad distribuida**: Los sistemas distribuidos son inherentemente más complejos.
2. **Consistencia de datos**: Mantener la consistencia entre servicios es desafiante.
3. **Observabilidad**: Se requieren herramientas especializadas para monitorizar y depurar.
4. **Latencia de red**: La comunicación entre servicios introduce latencia.
5. **Transacciones distribuidas**: Difíciles de implementar correctamente.

Python es particularmente adecuado para microservicios gracias a su desarrollo rápido, legibilidad y rico ecosistema de frameworks y bibliotecas. Con las herramientas y patrones adecuados, Python permite construir sistemas de microservicios escalables, resilientes y mantenibles.

## Ejercicios Prácticos

1. **Implementar un sistema básico de microservicios**:
   - Crear un servicio de usuarios con operaciones CRUD.
   - Crear un servicio de productos con operaciones CRUD.
   - Implementar un API Gateway que enrute las solicitudes a los servicios correspondientes.
   - Utilizar Docker y Docker Compose para ejecutar los servicios localmente.

2. **Añadir comunicación asíncrona**:
   - Implementar un broker de mensajes (RabbitMQ o Kafka).
   - Modificar los servicios para publicar eventos cuando ocurren cambios.
   - Crear un servicio de notificaciones que consuma eventos y envíe notificaciones.

3. **Implementar patrones de resiliencia**:
   - Añadir Circuit Breaker a las llamadas entre servicios.
   - Implementar reintentos con backoff exponencial.
   - Añadir timeouts a todas las operaciones externas.

4. **Añadir observabilidad**:
   - Configurar logging estructurado con correlationID.
   - Implementar métricas con Prometheus.
   - Añadir trazabilidad con Jaeger.
   - Crear dashboards en Grafana para visualizar métricas.

5. **Desplegar en Kubernetes**:
   - Crear manifiestos de Kubernetes para todos los servicios.
   - Configurar Ingress para el API Gateway.
   - Implementar health checks y readiness probes.
   - Configurar auto-scaling basado en métricas.

## Referencias

1. Sam Newman, "Building Microservices"
2. Chris Richardson, "Microservices Patterns"
3. Documentación oficial de Python: https://docs.python.org/
4. FastAPI: https://fastapi.tiangolo.com/
5. Kubernetes: https://kubernetes.io/docs/
6. Docker: https://docs.docker.com/
7. gRPC: https://grpc.io/docs/
8. Prometheus: https://prometheus.io/docs/
9. Jaeger: https://www.jaegertracing.io/docs/