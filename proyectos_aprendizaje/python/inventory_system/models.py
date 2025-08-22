#!/usr/bin/env python3
"""
Modelos para el Sistema de Gestión de Inventario

Este módulo define las clases de modelo para el sistema de gestión de inventario,
implementando el patrón de diseño Data Access Object (DAO) para separar la lógica
de negocio de la persistencia de datos.
"""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum, auto
from typing import Dict, List, Optional, Union
import json
import os
import uuid


class ProductCategory(Enum):
    """Enumeración de categorías de productos."""
    ELECTRONICS = auto()
    CLOTHING = auto()
    FOOD = auto()
    BOOKS = auto()
    TOYS = auto()
    HOME = auto()
    OFFICE = auto()
    OTHER = auto()
    
    def __str__(self) -> str:
        return self.name.capitalize()


class InventoryStatus(Enum):
    """Enumeración de estados de inventario."""
    IN_STOCK = auto()
    LOW_STOCK = auto()
    OUT_OF_STOCK = auto()
    DISCONTINUED = auto()
    
    def __str__(self) -> str:
        return self.name.replace('_', ' ').capitalize()


@dataclass
class Product:
    """Clase que representa un producto en el inventario."""
    name: str
    description: str
    category: ProductCategory
    price: float
    stock_quantity: int
    sku: str = field(default_factory=lambda: str(uuid.uuid4())[:8].upper())
    created_at: datetime = field(default_factory=datetime.now)
    updated_at: datetime = field(default_factory=datetime.now)
    
    @property
    def status(self) -> InventoryStatus:
        """Determina el estado del inventario basado en la cantidad disponible."""
        if self.stock_quantity <= 0:
            return InventoryStatus.OUT_OF_STOCK
        elif self.stock_quantity < 5:
            return InventoryStatus.LOW_STOCK
        else:
            return InventoryStatus.IN_STOCK
    
    def to_dict(self) -> Dict:
        """Convierte el objeto a un diccionario para serialización."""
        return {
            'sku': self.sku,
            'name': self.name,
            'description': self.description,
            'category': self.category.name,
            'price': self.price,
            'stock_quantity': self.stock_quantity,
            'created_at': self.created_at.isoformat(),
            'updated_at': self.updated_at.isoformat()
        }
    
    @classmethod
    def from_dict(cls, data: Dict) -> 'Product':
        """Crea un objeto Product a partir de un diccionario."""
        return cls(
            name=data['name'],
            description=data['description'],
            category=ProductCategory[data['category']],
            price=data['price'],
            stock_quantity=data['stock_quantity'],
            sku=data['sku'],
            created_at=datetime.fromisoformat(data['created_at']),
            updated_at=datetime.fromisoformat(data['updated_at'])
        )


@dataclass
class Transaction:
    """Clase que representa una transacción de inventario."""
    product_sku: str
    quantity: int
    transaction_type: str  # 'purchase', 'sale', 'adjustment'
    transaction_id: str = field(default_factory=lambda: str(uuid.uuid4()))
    timestamp: datetime = field(default_factory=datetime.now)
    notes: str = ""
    
    def to_dict(self) -> Dict:
        """Convierte el objeto a un diccionario para serialización."""
        return {
            'transaction_id': self.transaction_id,
            'product_sku': self.product_sku,
            'quantity': self.quantity,
            'transaction_type': self.transaction_type,
            'timestamp': self.timestamp.isoformat(),
            'notes': self.notes
        }
    
    @classmethod
    def from_dict(cls, data: Dict) -> 'Transaction':
        """Crea un objeto Transaction a partir de un diccionario."""
        return cls(
            product_sku=data['product_sku'],
            quantity=data['quantity'],
            transaction_type=data['transaction_type'],
            transaction_id=data['transaction_id'],
            timestamp=datetime.fromisoformat(data['timestamp']),
            notes=data.get('notes', "")
        )


class DAO(ABC):
    """Interfaz para el patrón Data Access Object."""
    
    @abstractmethod
    def save(self, obj: Union[Product, Transaction]) -> bool:
        """Guarda un objeto en el almacenamiento."""
        pass
    
    @abstractmethod
    def delete(self, id_value: str) -> bool:
        """Elimina un objeto del almacenamiento."""
        pass
    
    @abstractmethod
    def get(self, id_value: str) -> Optional[Union[Product, Transaction]]:
        """Obtiene un objeto del almacenamiento por su ID."""
        pass
    
    @abstractmethod
    def get_all(self) -> List[Union[Product, Transaction]]:
        """Obtiene todos los objetos del almacenamiento."""
        pass


class ProductDAO(DAO):
    """Implementación de DAO para productos."""
    
    def __init__(self, file_path: str = "products.json"):
        self.file_path = file_path
        self.products: Dict[str, Product] = {}
        self._load_data()
    
    def _load_data(self) -> None:
        """Carga los datos desde el archivo JSON."""
        if os.path.exists(self.file_path):
            try:
                with open(self.file_path, 'r') as f:
                    data = json.load(f)
                    for product_data in data:
                        product = Product.from_dict(product_data)
                        self.products[product.sku] = product
            except (json.JSONDecodeError, KeyError) as e:
                print(f"Error al cargar datos de productos: {e}")
    
    def _save_data(self) -> None:
        """Guarda los datos en el archivo JSON."""
        try:
            with open(self.file_path, 'w') as f:
                json.dump([p.to_dict() for p in self.products.values()], f, indent=2)
        except Exception as e:
            print(f"Error al guardar datos de productos: {e}")
    
    def save(self, product: Product) -> bool:
        """Guarda un producto en el almacenamiento."""
        if product.sku in self.products:
            product.updated_at = datetime.now()
        self.products[product.sku] = product
        self._save_data()
        return True
    
    def delete(self, sku: str) -> bool:
        """Elimina un producto del almacenamiento."""
        if sku in self.products:
            del self.products[sku]
            self._save_data()
            return True
        return False
    
    def get(self, sku: str) -> Optional[Product]:
        """Obtiene un producto del almacenamiento por su SKU."""
        return self.products.get(sku)
    
    def get_all(self) -> List[Product]:
        """Obtiene todos los productos del almacenamiento."""
        return list(self.products.values())
    
    def search(self, query: str) -> List[Product]:
        """Busca productos por nombre o descripción."""
        query = query.lower()
        return [p for p in self.products.values() 
                if query in p.name.lower() or query in p.description.lower()]
    
    def filter_by_category(self, category: ProductCategory) -> List[Product]:
        """Filtra productos por categoría."""
        return [p for p in self.products.values() if p.category == category]
    
    def filter_by_status(self, status: InventoryStatus) -> List[Product]:
        """Filtra productos por estado de inventario."""
        return [p for p in self.products.values() if p.status == status]


class TransactionDAO(DAO):
    """Implementación de DAO para transacciones."""
    
    def __init__(self, file_path: str = "transactions.json"):
        self.file_path = file_path
        self.transactions: Dict[str, Transaction] = {}
        self._load_data()
    
    def _load_data(self) -> None:
        """Carga los datos desde el archivo JSON."""
        if os.path.exists(self.file_path):
            try:
                with open(self.file_path, 'r') as f:
                    data = json.load(f)
                    for transaction_data in data:
                        transaction = Transaction.from_dict(transaction_data)
                        self.transactions[transaction.transaction_id] = transaction
            except (json.JSONDecodeError, KeyError) as e:
                print(f"Error al cargar datos de transacciones: {e}")
    
    def _save_data(self) -> None:
        """Guarda los datos en el archivo JSON."""
        try:
            with open(self.file_path, 'w') as f:
                json.dump([t.to_dict() for t in self.transactions.values()], f, indent=2)
        except Exception as e:
            print(f"Error al guardar datos de transacciones: {e}")
    
    def save(self, transaction: Transaction) -> bool:
        """Guarda una transacción en el almacenamiento."""
        self.transactions[transaction.transaction_id] = transaction
        self._save_data()
        return True
    
    def delete(self, transaction_id: str) -> bool:
        """Elimina una transacción del almacenamiento."""
        if transaction_id in self.transactions:
            del self.transactions[transaction_id]
            self._save_data()
            return True
        return False
    
    def get(self, transaction_id: str) -> Optional[Transaction]:
        """Obtiene una transacción del almacenamiento por su ID."""
        return self.transactions.get(transaction_id)
    
    def get_all(self) -> List[Transaction]:
        """Obtiene todas las transacciones del almacenamiento."""
        return list(self.transactions.values())
    
    def get_by_product(self, product_sku: str) -> List[Transaction]:
        """Obtiene todas las transacciones para un producto específico."""
        return [t for t in self.transactions.values() if t.product_sku == product_sku]
    
    def get_by_type(self, transaction_type: str) -> List[Transaction]:
        """Obtiene todas las transacciones de un tipo específico."""
        return [t for t in self.transactions.values() if t.transaction_type == transaction_type]