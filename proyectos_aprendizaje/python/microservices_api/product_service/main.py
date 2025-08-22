from fastapi import FastAPI, HTTPException, Depends, status
from pydantic import BaseModel, Field
from typing import List, Optional
import uuid
import json
import os
from datetime import datetime
import uvicorn

# Modelos
class ProductBase(BaseModel):
    name: str
    description: Optional[str] = None
    price: float = Field(gt=0)
    category: str
    stock: int = Field(ge=0)

class ProductCreate(ProductBase):
    pass

class ProductUpdate(BaseModel):
    name: Optional[str] = None
    description: Optional[str] = None
    price: Optional[float] = Field(None, gt=0)
    category: Optional[str] = None
    stock: Optional[int] = Field(None, ge=0)

class Product(ProductBase):
    id: str
    created_at: datetime
    updated_at: datetime

    class Config:
        orm_mode = True

# Simulación de base de datos
class ProductDB:
    def __init__(self, file_path="products.json"):
        self.file_path = file_path
        self.ensure_file_exists()

    def ensure_file_exists(self):
        if not os.path.exists(self.file_path):
            with open(self.file_path, "w") as f:
                json.dump([], f)

    def read_products(self):
        with open(self.file_path, "r") as f:
            try:
                return json.load(f)
            except json.JSONDecodeError:
                return []

    def write_products(self, products):
        with open(self.file_path, "w") as f:
            json.dump(products, f, default=str)

    def get_all(self):
        return self.read_products()

    def get_by_id(self, product_id):
        products = self.read_products()
        for product in products:
            if product["id"] == product_id:
                return product
        return None

    def create(self, product_data):
        products = self.read_products()
        new_product = product_data.dict()
        new_product["id"] = str(uuid.uuid4())
        new_product["created_at"] = datetime.now()
        new_product["updated_at"] = datetime.now()
        products.append(new_product)
        self.write_products(products)
        return new_product

    def update(self, product_id, product_data):
        products = self.read_products()
        for i, product in enumerate(products):
            if product["id"] == product_id:
                update_data = product_data.dict(exclude_unset=True)
                products[i].update(update_data)
                products[i]["updated_at"] = datetime.now()
                self.write_products(products)
                return products[i]
        return None

    def delete(self, product_id):
        products = self.read_products()
        for i, product in enumerate(products):
            if product["id"] == product_id:
                deleted_product = products.pop(i)
                self.write_products(products)
                return deleted_product
        return None

    def get_by_category(self, category):
        products = self.read_products()
        return [product for product in products if product["category"] == category]

# Inicializar la aplicación
app = FastAPI(title="Product Service", description="Servicio de gestión de productos para la API de microservicios")

# Dependencia para obtener la base de datos
def get_db():
    return ProductDB()

# Rutas
@app.get("/products/", response_model=List[Product])
async def read_products(skip: int = 0, limit: int = 100, db: ProductDB = Depends(get_db)):
    products = db.get_all()
    return products[skip : skip + limit]

@app.get("/products/{product_id}", response_model=Product)
async def read_product(product_id: str, db: ProductDB = Depends(get_db)):
    product = db.get_by_id(product_id)
    if product is None:
        raise HTTPException(status_code=404, detail="Product not found")
    return product

@app.post("/products/", response_model=Product, status_code=status.HTTP_201_CREATED)
async def create_product(product: ProductCreate, db: ProductDB = Depends(get_db)):
    return db.create(product)

@app.put("/products/{product_id}", response_model=Product)
async def update_product(product_id: str, product: ProductUpdate, db: ProductDB = Depends(get_db)):
    updated_product = db.update(product_id, product)
    if updated_product is None:
        raise HTTPException(status_code=404, detail="Product not found")
    return updated_product

@app.delete("/products/{product_id}", response_model=Product)
async def delete_product(product_id: str, db: ProductDB = Depends(get_db)):
    deleted_product = db.delete(product_id)
    if deleted_product is None:
        raise HTTPException(status_code=404, detail="Product not found")
    return deleted_product

@app.get("/products/category/{category}", response_model=List[Product])
async def read_products_by_category(category: str, db: ProductDB = Depends(get_db)):
    return db.get_by_category(category)

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "product"}

@app.get("/")
async def root():
    return {"message": "Product Service API"}

# Para ejecutar directamente este archivo
if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)