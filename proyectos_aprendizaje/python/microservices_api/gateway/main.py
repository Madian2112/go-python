from fastapi import FastAPI, HTTPException, Depends, Request, status
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Dict, Any, Optional, List
import httpx
import uvicorn
import os
from datetime import datetime

# Configuración de servicios
SERVICES = {
    "auth": "http://localhost:8000",
    "product": "http://localhost:8001",
    "order": "http://localhost:8002",
}

# Inicializar la aplicación
app = FastAPI(title="API Gateway", description="Gateway para la API de microservicios")

# Configurar CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # En producción, especificar orígenes permitidos
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Cliente HTTP asíncrono
async_client = httpx.AsyncClient()

# Modelos
class ServiceStatus(BaseModel):
    service: str
    status: str
    url: str
    latency_ms: float

# Middleware para logging
@app.middleware("http")
async def log_requests(request: Request, call_next):
    start_time = datetime.now()
    response = await call_next(request)
    process_time = (datetime.now() - start_time).total_seconds() * 1000
    print(f"[{datetime.now()}] {request.method} {request.url.path} - {response.status_code} - {process_time:.2f}ms")
    return response

# Funciones de utilidad
async def forward_request(service: str, path: str, method: str, data: Any = None, params: Dict = None, headers: Dict = None):
    if service not in SERVICES:
        raise HTTPException(status_code=404, detail=f"Service '{service}' not found")
    
    url = f"{SERVICES[service]}{path}"
    
    try:
        if method.lower() == "get":
            response = await async_client.get(url, params=params, headers=headers)
        elif method.lower() == "post":
            response = await async_client.post(url, json=data, params=params, headers=headers)
        elif method.lower() == "put":
            response = await async_client.put(url, json=data, params=params, headers=headers)
        elif method.lower() == "delete":
            response = await async_client.delete(url, params=params, headers=headers)
        else:
            raise HTTPException(status_code=405, detail=f"Method '{method}' not allowed")
        
        return JSONResponse(
            content=response.json(),
            status_code=response.status_code,
            headers=dict(response.headers)
        )
    except httpx.RequestError as e:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail=f"Error communicating with service: {str(e)}"
        )

# Rutas
@app.get("/")
async def root():
    return {"message": "API Gateway for Microservices"}

@app.get("/health", response_model=List[ServiceStatus])
async def health_check():
    results = []
    for service_name, service_url in SERVICES.items():
        try:
            start_time = datetime.now()
            response = await async_client.get(f"{service_url}/health")
            latency = (datetime.now() - start_time).total_seconds() * 1000
            
            if response.status_code == 200:
                results.append({
                    "service": service_name,
                    "status": "healthy",
                    "url": service_url,
                    "latency_ms": latency
                })
            else:
                results.append({
                    "service": service_name,
                    "status": "unhealthy",
                    "url": service_url,
                    "latency_ms": latency
                })
        except Exception as e:
            results.append({
                "service": service_name,
                "status": "unavailable",
                "url": service_url,
                "latency_ms": 0
            })
    
    return results

# Rutas para el servicio de autenticación
@app.post("/api/auth/token")
async def login(request: Request):
    data = await request.json()
    return await forward_request("auth", "/token", "post", data)

@app.get("/api/auth/users/me")
async def get_current_user(request: Request):
    headers = dict(request.headers)
    return await forward_request("auth", "/users/me/", "get", headers=headers)

# Rutas para el servicio de productos
@app.get("/api/products")
async def get_products(request: Request):
    params = dict(request.query_params)
    return await forward_request("product", "/products/", "get", params=params)

@app.get("/api/products/{product_id}")
async def get_product(product_id: str, request: Request):
    return await forward_request("product", f"/products/{product_id}", "get")

@app.post("/api/products")
async def create_product(request: Request):
    data = await request.json()
    headers = dict(request.headers)
    return await forward_request("product", "/products/", "post", data, headers=headers)

@app.put("/api/products/{product_id}")
async def update_product(product_id: str, request: Request):
    data = await request.json()
    headers = dict(request.headers)
    return await forward_request("product", f"/products/{product_id}", "put", data, headers=headers)

@app.delete("/api/products/{product_id}")
async def delete_product(product_id: str, request: Request):
    headers = dict(request.headers)
    return await forward_request("product", f"/products/{product_id}", "delete", headers=headers)

@app.get("/api/products/category/{category}")
async def get_products_by_category(category: str, request: Request):
    return await forward_request("product", f"/products/category/{category}", "get")

# Rutas para el servicio de órdenes
@app.get("/api/orders")
async def get_orders(request: Request):
    params = dict(request.query_params)
    headers = dict(request.headers)
    return await forward_request("order", "/orders/", "get", params=params, headers=headers)

@app.get("/api/orders/{order_id}")
async def get_order(order_id: str, request: Request):
    headers = dict(request.headers)
    return await forward_request("order", f"/orders/{order_id}", "get", headers=headers)

@app.post("/api/orders")
async def create_order(request: Request):
    data = await request.json()
    headers = dict(request.headers)
    return await forward_request("order", "/orders/", "post", data, headers=headers)

@app.put("/api/orders/{order_id}")
async def update_order(order_id: str, request: Request):
    data = await request.json()
    headers = dict(request.headers)
    return await forward_request("order", f"/orders/{order_id}", "put", data, headers=headers)

@app.delete("/api/orders/{order_id}")
async def delete_order(order_id: str, request: Request):
    headers = dict(request.headers)
    return await forward_request("order", f"/orders/{order_id}", "delete", headers=headers)

@app.get("/api/orders/customer/{customer_id}")
async def get_orders_by_customer(customer_id: str, request: Request):
    headers = dict(request.headers)
    return await forward_request("order", f"/orders/customer/{customer_id}", "get", headers=headers)

# Para ejecutar directamente este archivo
if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8080, reload=True)