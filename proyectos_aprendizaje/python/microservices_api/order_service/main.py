from fastapi import FastAPI, HTTPException, Depends, status
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import uuid
import json
import os
from datetime import datetime
import requests
import uvicorn

# Configuración
PRODUCT_SERVICE_URL = "http://localhost:8001"

# Modelos
class OrderItemCreate(BaseModel):
    product_id: str
    quantity: int = Field(gt=0)

class OrderItem(OrderItemCreate):
    product_name: str
    product_price: float
    subtotal: float

class OrderCreate(BaseModel):
    customer_id: str
    items: List[OrderItemCreate]
    shipping_address: str

class OrderStatus(str):
    PENDING = "PENDING"
    PROCESSING = "PROCESSING"
    SHIPPED = "SHIPPED"
    DELIVERED = "DELIVERED"
    CANCELLED = "CANCELLED"

class Order(BaseModel):
    id: str
    customer_id: str
    items: List[OrderItem]
    total: float
    status: str
    shipping_address: str
    created_at: datetime
    updated_at: datetime

    class Config:
        orm_mode = True

class OrderUpdate(BaseModel):
    status: Optional[str] = None
    shipping_address: Optional[str] = None

# Simulación de base de datos
class OrderDB:
    def __init__(self, file_path="orders.json"):
        self.file_path = file_path
        self.ensure_file_exists()

    def ensure_file_exists(self):
        if not os.path.exists(self.file_path):
            with open(self.file_path, "w") as f:
                json.dump([], f)

    def read_orders(self):
        with open(self.file_path, "r") as f:
            try:
                return json.load(f)
            except json.JSONDecodeError:
                return []

    def write_orders(self, orders):
        with open(self.file_path, "w") as f:
            json.dump(orders, f, default=str)

    def get_all(self):
        return self.read_orders()

    def get_by_id(self, order_id):
        orders = self.read_orders()
        for order in orders:
            if order["id"] == order_id:
                return order
        return None

    def create(self, order_data, processed_items):
        orders = self.read_orders()
        new_order = {
            "id": str(uuid.uuid4()),
            "customer_id": order_data.customer_id,
            "items": processed_items,
            "total": sum(item["subtotal"] for item in processed_items),
            "status": OrderStatus.PENDING,
            "shipping_address": order_data.shipping_address,
            "created_at": datetime.now(),
            "updated_at": datetime.now()
        }
        orders.append(new_order)
        self.write_orders(orders)
        return new_order

    def update(self, order_id, order_data):
        orders = self.read_orders()
        for i, order in enumerate(orders):
            if order["id"] == order_id:
                update_data = order_data.dict(exclude_unset=True)
                orders[i].update(update_data)
                orders[i]["updated_at"] = datetime.now()
                self.write_orders(orders)
                return orders[i]
        return None

    def delete(self, order_id):
        orders = self.read_orders()
        for i, order in enumerate(orders):
            if order["id"] == order_id:
                deleted_order = orders.pop(i)
                self.write_orders(orders)
                return deleted_order
        return None

    def get_by_customer(self, customer_id):
        orders = self.read_orders()
        return [order for order in orders if order["customer_id"] == customer_id]

# Inicializar la aplicación
app = FastAPI(title="Order Service", description="Servicio de gestión de órdenes para la API de microservicios")

# Dependencia para obtener la base de datos
def get_db():
    return OrderDB()

# Funciones de utilidad
async def get_product_details(product_id: str):
    try:
        response = requests.get(f"{PRODUCT_SERVICE_URL}/products/{product_id}")
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail=f"Error communicating with product service: {str(e)}"
        )

async def process_order_items(items: List[OrderItemCreate]):
    processed_items = []
    for item in items:
        product = await get_product_details(item.product_id)
        processed_item = {
            "product_id": item.product_id,
            "product_name": product["name"],
            "product_price": product["price"],
            "quantity": item.quantity,
            "subtotal": product["price"] * item.quantity
        }
        processed_items.append(processed_item)
    return processed_items

# Rutas
@app.get("/orders/", response_model=List[Order])
async def read_orders(skip: int = 0, limit: int = 100, db: OrderDB = Depends(get_db)):
    orders = db.get_all()
    return orders[skip : skip + limit]

@app.get("/orders/{order_id}", response_model=Order)
async def read_order(order_id: str, db: OrderDB = Depends(get_db)):
    order = db.get_by_id(order_id)
    if order is None:
        raise HTTPException(status_code=404, detail="Order not found")
    return order

@app.post("/orders/", response_model=Order, status_code=status.HTTP_201_CREATED)
async def create_order(order: OrderCreate, db: OrderDB = Depends(get_db)):
    try:
        processed_items = await process_order_items(order.items)
        return db.create(order, processed_items)
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Error creating order: {str(e)}"
        )

@app.put("/orders/{order_id}", response_model=Order)
async def update_order(order_id: str, order: OrderUpdate, db: OrderDB = Depends(get_db)):
    updated_order = db.update(order_id, order)
    if updated_order is None:
        raise HTTPException(status_code=404, detail="Order not found")
    return updated_order

@app.delete("/orders/{order_id}", response_model=Order)
async def delete_order(order_id: str, db: OrderDB = Depends(get_db)):
    deleted_order = db.delete(order_id)
    if deleted_order is None:
        raise HTTPException(status_code=404, detail="Order not found")
    return deleted_order

@app.get("/orders/customer/{customer_id}", response_model=List[Order])
async def read_orders_by_customer(customer_id: str, db: OrderDB = Depends(get_db)):
    return db.get_by_customer(customer_id)

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "order"}

@app.get("/")
async def root():
    return {"message": "Order Service API"}

# Para ejecutar directamente este archivo
if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8002, reload=True)