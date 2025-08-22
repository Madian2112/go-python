#!/usr/bin/env python3
"""
Servicios para el Sistema de Gestión de Inventario

Este módulo implementa el patrón de diseño Facade para proporcionar una interfaz
simplificada al sistema de gestión de inventario, ocultando la complejidad de las
interacciones entre los diferentes componentes.
"""

from datetime import datetime
from typing import Dict, List, Optional, Tuple

from models import (
    Product, Transaction, ProductCategory, InventoryStatus,
    ProductDAO, TransactionDAO
)


class InventoryService:
    """Facade para el sistema de gestión de inventario."""
    
    def __init__(self, product_dao: ProductDAO = None, transaction_dao: TransactionDAO = None):
        """Inicializa el servicio con los DAOs especificados o crea nuevos."""
        self.product_dao = product_dao or ProductDAO()
        self.transaction_dao = transaction_dao or TransactionDAO()
    
    # Métodos para gestión de productos
    
    def add_product(self, name: str, description: str, category: str, 
                   price: float, stock_quantity: int) -> Product:
        """Añade un nuevo producto al inventario."""
        try:
            category_enum = ProductCategory[category.upper()]
        except KeyError:
            raise ValueError(f"Categoría inválida: {category}. Opciones válidas: {', '.join([c.name for c in ProductCategory])}")
        
        if price <= 0:
            raise ValueError("El precio debe ser mayor que cero")
        
        if stock_quantity < 0:
            raise ValueError("La cantidad en stock no puede ser negativa")
        
        product = Product(
            name=name,
            description=description,
            category=category_enum,
            price=price,
            stock_quantity=stock_quantity
        )
        
        self.product_dao.save(product)
        
        # Registrar transacción de compra inicial
        if stock_quantity > 0:
            self._record_transaction(
                product.sku, 
                stock_quantity, 
                "purchase", 
                f"Inventario inicial para {product.name}"
            )
        
        return product
    
    def update_product(self, sku: str, **kwargs) -> Optional[Product]:
        """Actualiza un producto existente."""
        product = self.product_dao.get(sku)
        if not product:
            return None
        
        # Actualizar los campos proporcionados
        if 'name' in kwargs:
            product.name = kwargs['name']
        
        if 'description' in kwargs:
            product.description = kwargs['description']
        
        if 'category' in kwargs:
            try:
                product.category = ProductCategory[kwargs['category'].upper()]
            except KeyError:
                raise ValueError(f"Categoría inválida: {kwargs['category']}")
        
        if 'price' in kwargs:
            if kwargs['price'] <= 0:
                raise ValueError("El precio debe ser mayor que cero")
            product.price = kwargs['price']
        
        # No actualizamos stock_quantity directamente, se debe usar add_stock o remove_stock
        
        product.updated_at = datetime.now()
        self.product_dao.save(product)
        return product
    
    def get_product(self, sku: str) -> Optional[Product]:
        """Obtiene un producto por su SKU."""
        return self.product_dao.get(sku)
    
    def delete_product(self, sku: str) -> bool:
        """Elimina un producto del inventario."""
        product = self.product_dao.get(sku)
        if not product:
            return False
        
        # Registrar transacción de ajuste si hay stock
        if product.stock_quantity > 0:
            self._record_transaction(
                sku, 
                -product.stock_quantity, 
                "adjustment", 
                f"Eliminación del producto {product.name}"
            )
        
        return self.product_dao.delete(sku)
    
    def search_products(self, query: str) -> List[Product]:
        """Busca productos por nombre o descripción."""
        return self.product_dao.search(query)
    
    def get_all_products(self) -> List[Product]:
        """Obtiene todos los productos."""
        return self.product_dao.get_all()
    
    def get_products_by_category(self, category: str) -> List[Product]:
        """Obtiene productos por categoría."""
        try:
            category_enum = ProductCategory[category.upper()]
            return self.product_dao.filter_by_category(category_enum)
        except KeyError:
            raise ValueError(f"Categoría inválida: {category}")
    
    def get_products_by_status(self, status: str) -> List[Product]:
        """Obtiene productos por estado de inventario."""
        try:
            status_enum = InventoryStatus[status.upper()]
            return self.product_dao.filter_by_status(status_enum)
        except KeyError:
            raise ValueError(f"Estado inválido: {status}")
    
    # Métodos para gestión de stock
    
    def add_stock(self, sku: str, quantity: int, notes: str = "") -> Tuple[bool, str]:
        """Añade stock a un producto existente."""
        if quantity <= 0:
            return False, "La cantidad debe ser mayor que cero"
        
        product = self.product_dao.get(sku)
        if not product:
            return False, f"Producto con SKU {sku} no encontrado"
        
        product.stock_quantity += quantity
        product.updated_at = datetime.now()
        self.product_dao.save(product)
        
        # Registrar transacción
        self._record_transaction(sku, quantity, "purchase", notes)
        
        return True, f"Stock actualizado. Nuevo stock: {product.stock_quantity}"
    
    def remove_stock(self, sku: str, quantity: int, notes: str = "") -> Tuple[bool, str]:
        """Reduce el stock de un producto existente."""
        if quantity <= 0:
            return False, "La cantidad debe ser mayor que cero"
        
        product = self.product_dao.get(sku)
        if not product:
            return False, f"Producto con SKU {sku} no encontrado"
        
        if product.stock_quantity < quantity:
            return False, f"Stock insuficiente. Disponible: {product.stock_quantity}"
        
        product.stock_quantity -= quantity
        product.updated_at = datetime.now()
        self.product_dao.save(product)
        
        # Registrar transacción
        self._record_transaction(sku, -quantity, "sale", notes)
        
        return True, f"Stock actualizado. Nuevo stock: {product.stock_quantity}"
    
    def adjust_stock(self, sku: str, new_quantity: int, notes: str = "") -> Tuple[bool, str]:
        """Ajusta el stock de un producto a una cantidad específica."""
        if new_quantity < 0:
            return False, "La cantidad no puede ser negativa"
        
        product = self.product_dao.get(sku)
        if not product:
            return False, f"Producto con SKU {sku} no encontrado"
        
        difference = new_quantity - product.stock_quantity
        product.stock_quantity = new_quantity
        product.updated_at = datetime.now()
        self.product_dao.save(product)
        
        # Registrar transacción de ajuste
        self._record_transaction(sku, difference, "adjustment", notes)
        
        return True, f"Stock ajustado a {new_quantity}"
    
    # Métodos para gestión de transacciones
    
    def _record_transaction(self, product_sku: str, quantity: int, 
                           transaction_type: str, notes: str = "") -> Transaction:
        """Registra una transacción de inventario."""
        transaction = Transaction(
            product_sku=product_sku,
            quantity=quantity,
            transaction_type=transaction_type,
            notes=notes
        )
        self.transaction_dao.save(transaction)
        return transaction
    
    def get_transaction(self, transaction_id: str) -> Optional[Transaction]:
        """Obtiene una transacción por su ID."""
        return self.transaction_dao.get(transaction_id)
    
    def get_all_transactions(self) -> List[Transaction]:
        """Obtiene todas las transacciones."""
        return self.transaction_dao.get_all()
    
    def get_product_transactions(self, sku: str) -> List[Transaction]:
        """Obtiene todas las transacciones para un producto específico."""
        return self.transaction_dao.get_by_product(sku)
    
    def get_transactions_by_type(self, transaction_type: str) -> List[Transaction]:
        """Obtiene todas las transacciones de un tipo específico."""
        valid_types = ["purchase", "sale", "adjustment"]
        if transaction_type not in valid_types:
            raise ValueError(f"Tipo de transacción inválido. Opciones válidas: {', '.join(valid_types)}")
        
        return self.transaction_dao.get_by_type(transaction_type)
    
    # Métodos para reportes
    
    def get_low_stock_products(self, threshold: int = 5) -> List[Product]:
        """Obtiene productos con stock bajo."""
        return [p for p in self.product_dao.get_all() if p.stock_quantity <= threshold]
    
    def get_out_of_stock_products(self) -> List[Product]:
        """Obtiene productos sin stock."""
        return [p for p in self.product_dao.get_all() if p.stock_quantity <= 0]
    
    def get_inventory_value(self) -> Dict[str, float]:
        """Calcula el valor total del inventario."""
        products = self.product_dao.get_all()
        total_value = sum(p.price * p.stock_quantity for p in products)
        by_category = {}
        
        for category in ProductCategory:
            category_products = [p for p in products if p.category == category]
            category_value = sum(p.price * p.stock_quantity for p in category_products)
            by_category[category.name] = category_value
        
        return {
            "total": total_value,
            "by_category": by_category
        }
    
    def get_transaction_summary(self) -> Dict[str, int]:
        """Genera un resumen de transacciones por tipo."""
        transactions = self.transaction_dao.get_all()
        summary = {
            "purchase": 0,
            "sale": 0,
            "adjustment": 0
        }
        
        for transaction in transactions:
            if transaction.transaction_type in summary:
                summary[transaction.transaction_type] += 1
        
        return summary