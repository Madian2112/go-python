#!/usr/bin/env python3
"""
Interfaz de Línea de Comandos para el Sistema de Gestión de Inventario

Este módulo implementa una interfaz de línea de comandos para interactuar con
el sistema de gestión de inventario, utilizando el patrón de diseño Command
para encapsular las operaciones como objetos.
"""

import argparse
import os
import sys
from abc import ABC, abstractmethod
from datetime import datetime
from typing import Dict, List, Optional, Any

from models import Product, ProductCategory, InventoryStatus
from services import InventoryService


class Command(ABC):
    """Interfaz para el patrón Command."""
    
    @abstractmethod
    def execute(self) -> None:
        """Ejecuta el comando."""
        pass


class AddProductCommand(Command):
    """Comando para añadir un nuevo producto."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        try:
            product = self.service.add_product(
                name=self.args.name,
                description=self.args.description,
                category=self.args.category,
                price=float(self.args.price),
                stock_quantity=int(self.args.quantity)
            )
            print(f"\nProducto añadido correctamente:")
            print_product_details(product)
        except ValueError as e:
            print(f"\nError: {e}")


class UpdateProductCommand(Command):
    """Comando para actualizar un producto existente."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        updates = {}
        if hasattr(self.args, 'name') and self.args.name is not None:
            updates['name'] = self.args.name
        
        if hasattr(self.args, 'description') and self.args.description is not None:
            updates['description'] = self.args.description
        
        if hasattr(self.args, 'category') and self.args.category is not None:
            updates['category'] = self.args.category
        
        if hasattr(self.args, 'price') and self.args.price is not None:
            updates['price'] = float(self.args.price)
        
        try:
            product = self.service.update_product(self.args.sku, **updates)
            if product:
                print(f"\nProducto actualizado correctamente:")
                print_product_details(product)
            else:
                print(f"\nError: Producto con SKU {self.args.sku} no encontrado")
        except ValueError as e:
            print(f"\nError: {e}")


class DeleteProductCommand(Command):
    """Comando para eliminar un producto."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        if self.service.delete_product(self.args.sku):
            print(f"\nProducto con SKU {self.args.sku} eliminado correctamente")
        else:
            print(f"\nError: Producto con SKU {self.args.sku} no encontrado")


class ListProductsCommand(Command):
    """Comando para listar productos."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        products = []
        
        if hasattr(self.args, 'category') and self.args.category:
            try:
                products = self.service.get_products_by_category(self.args.category)
            except ValueError as e:
                print(f"\nError: {e}")
                return
        elif hasattr(self.args, 'status') and self.args.status:
            try:
                products = self.service.get_products_by_status(self.args.status)
            except ValueError as e:
                print(f"\nError: {e}")
                return
        elif hasattr(self.args, 'query') and self.args.query:
            products = self.service.search_products(self.args.query)
        else:
            products = self.service.get_all_products()
        
        if not products:
            print("\nNo se encontraron productos")
            return
        
        print(f"\nSe encontraron {len(products)} productos:")
        print("-" * 80)
        for product in products:
            print_product_summary(product)
            print("-" * 80)


class ShowProductCommand(Command):
    """Comando para mostrar detalles de un producto."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        product = self.service.get_product(self.args.sku)
        if product:
            print("\nDetalles del producto:")
            print_product_details(product)
            
            # Mostrar historial de transacciones si se solicita
            if hasattr(self.args, 'transactions') and self.args.transactions:
                transactions = self.service.get_product_transactions(self.args.sku)
                if transactions:
                    print("\nHistorial de transacciones:")
                    print("-" * 80)
                    for transaction in transactions:
                        print_transaction(transaction)
                else:
                    print("\nNo hay transacciones registradas para este producto")
        else:
            print(f"\nError: Producto con SKU {self.args.sku} no encontrado")


class AddStockCommand(Command):
    """Comando para añadir stock a un producto."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        notes = self.args.notes if hasattr(self.args, 'notes') and self.args.notes else ""
        success, message = self.service.add_stock(
            self.args.sku, 
            int(self.args.quantity),
            notes
        )
        
        if success:
            print(f"\nStock añadido correctamente: {message}")
            product = self.service.get_product(self.args.sku)
            if product:
                print_product_summary(product)
        else:
            print(f"\nError: {message}")


class RemoveStockCommand(Command):
    """Comando para reducir stock de un producto."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        notes = self.args.notes if hasattr(self.args, 'notes') and self.args.notes else ""
        success, message = self.service.remove_stock(
            self.args.sku, 
            int(self.args.quantity),
            notes
        )
        
        if success:
            print(f"\nStock reducido correctamente: {message}")
            product = self.service.get_product(self.args.sku)
            if product:
                print_product_summary(product)
        else:
            print(f"\nError: {message}")


class AdjustStockCommand(Command):
    """Comando para ajustar stock de un producto."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        notes = self.args.notes if hasattr(self.args, 'notes') and self.args.notes else ""
        success, message = self.service.adjust_stock(
            self.args.sku, 
            int(self.args.quantity),
            notes
        )
        
        if success:
            print(f"\nStock ajustado correctamente: {message}")
            product = self.service.get_product(self.args.sku)
            if product:
                print_product_summary(product)
        else:
            print(f"\nError: {message}")


class ListTransactionsCommand(Command):
    """Comando para listar transacciones."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        transactions = []
        
        if hasattr(self.args, 'sku') and self.args.sku:
            transactions = self.service.get_product_transactions(self.args.sku)
        elif hasattr(self.args, 'type') and self.args.type:
            try:
                transactions = self.service.get_transactions_by_type(self.args.type)
            except ValueError as e:
                print(f"\nError: {e}")
                return
        else:
            transactions = self.service.get_all_transactions()
        
        if not transactions:
            print("\nNo se encontraron transacciones")
            return
        
        print(f"\nSe encontraron {len(transactions)} transacciones:")
        print("-" * 80)
        for transaction in transactions:
            print_transaction(transaction)
            print("-" * 80)


class ReportCommand(Command):
    """Comando para generar reportes."""
    
    def __init__(self, service: InventoryService, args: argparse.Namespace):
        self.service = service
        self.args = args
    
    def execute(self) -> None:
        if self.args.report_type == "low-stock":
            threshold = int(self.args.threshold) if hasattr(self.args, 'threshold') and self.args.threshold else 5
            products = self.service.get_low_stock_products(threshold)
            if products:
                print(f"\nProductos con stock bajo (menos de {threshold} unidades):")
                print("-" * 80)
                for product in products:
                    print_product_summary(product)
                    print("-" * 80)
            else:
                print(f"\nNo hay productos con stock bajo (menos de {threshold} unidades)")
        
        elif self.args.report_type == "out-of-stock":
            products = self.service.get_out_of_stock_products()
            if products:
                print("\nProductos sin stock:")
                print("-" * 80)
                for product in products:
                    print_product_summary(product)
                    print("-" * 80)
            else:
                print("\nNo hay productos sin stock")
        
        elif self.args.report_type == "inventory-value":
            value_data = self.service.get_inventory_value()
            print("\nValor del inventario:")
            print(f"Valor total: ${value_data['total']:.2f}")
            print("\nValor por categoría:")
            for category, value in value_data['by_category'].items():
                print(f"  {category}: ${value:.2f}")
        
        elif self.args.report_type == "transaction-summary":
            summary = self.service.get_transaction_summary()
            print("\nResumen de transacciones:")
            print(f"Compras: {summary['purchase']}")
            print(f"Ventas: {summary['sale']}")
            print(f"Ajustes: {summary['adjustment']}")
            print(f"Total: {sum(summary.values())}")


# Funciones auxiliares para imprimir información

def print_product_summary(product: Product) -> None:
    """Imprime un resumen de un producto."""
    status_colors = {
        InventoryStatus.IN_STOCK: "\033[92m",      # Verde
        InventoryStatus.LOW_STOCK: "\033[93m",     # Amarillo
        InventoryStatus.OUT_OF_STOCK: "\033[91m",  # Rojo
        InventoryStatus.DISCONTINUED: "\033[90m"   # Gris
    }
    
    color = status_colors.get(product.status, "\033[0m")
    reset = "\033[0m"
    
    print(f"SKU: {product.sku} | {product.name}")
    print(f"Categoría: {product.category} | Precio: ${product.price:.2f}")
    print(f"Stock: {color}{product.stock_quantity} ({product.status}){reset}")


def print_product_details(product: Product) -> None:
    """Imprime detalles completos de un producto."""
    status_colors = {
        InventoryStatus.IN_STOCK: "\033[92m",      # Verde
        InventoryStatus.LOW_STOCK: "\033[93m",     # Amarillo
        InventoryStatus.OUT_OF_STOCK: "\033[91m",  # Rojo
        InventoryStatus.DISCONTINUED: "\033[90m"   # Gris
    }
    
    color = status_colors.get(product.status, "\033[0m")
    reset = "\033[0m"
    
    print("-" * 80)
    print(f"SKU: {product.sku}")
    print(f"Nombre: {product.name}")
    print(f"Descripción: {product.description}")
    print(f"Categoría: {product.category}")
    print(f"Precio: ${product.price:.2f}")
    print(f"Stock: {color}{product.stock_quantity} ({product.status}){reset}")
    print(f"Creado: {product.created_at.strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"Actualizado: {product.updated_at.strftime('%Y-%m-%d %H:%M:%S')}")
    print("-" * 80)


def print_transaction(transaction: Any) -> None:
    """Imprime detalles de una transacción."""
    transaction_colors = {
        "purchase": "\033[92m",    # Verde
        "sale": "\033[91m",       # Rojo
        "adjustment": "\033[93m"  # Amarillo
    }
    
    color = transaction_colors.get(transaction.transaction_type, "\033[0m")
    reset = "\033[0m"
    
    quantity_str = f"+{transaction.quantity}" if transaction.quantity > 0 else str(transaction.quantity)
    
    print(f"ID: {transaction.transaction_id}")
    print(f"Producto: {transaction.product_sku}")
    print(f"Tipo: {color}{transaction.transaction_type.capitalize()}{reset}")
    print(f"Cantidad: {color}{quantity_str}{reset}")
    print(f"Fecha: {transaction.timestamp.strftime('%Y-%m-%d %H:%M:%S')}")
    if transaction.notes:
        print(f"Notas: {transaction.notes}")


class InventoryCLI:
    """Interfaz de línea de comandos para el sistema de gestión de inventario."""
    
    def __init__(self):
        self.service = InventoryService()
        self.parser = self._create_parser()
    
    def _create_parser(self) -> argparse.ArgumentParser:
        """Crea y configura el parser de argumentos."""
        parser = argparse.ArgumentParser(
            description="Sistema de Gestión de Inventario",
            formatter_class=argparse.RawDescriptionHelpFormatter
        )
        
        subparsers = parser.add_subparsers(dest="command", help="Comandos disponibles")
        
        # Comando: add-product
        add_product_parser = subparsers.add_parser("add-product", help="Añadir un nuevo producto")
        add_product_parser.add_argument("--name", "-n", required=True, help="Nombre del producto")
        add_product_parser.add_argument("--description", "-d", required=True, help="Descripción del producto")
        add_product_parser.add_argument("--category", "-c", required=True, 
                                      help=f"Categoría del producto. Opciones: {', '.join([c.name for c in ProductCategory])}")
        add_product_parser.add_argument("--price", "-p", required=True, help="Precio del producto")
        add_product_parser.add_argument("--quantity", "-q", required=True, help="Cantidad inicial en stock")
        
        # Comando: update-product
        update_product_parser = subparsers.add_parser("update-product", help="Actualizar un producto existente")
        update_product_parser.add_argument("--sku", "-s", required=True, help="SKU del producto")
        update_product_parser.add_argument("--name", "-n", help="Nuevo nombre del producto")
        update_product_parser.add_argument("--description", "-d", help="Nueva descripción del producto")
        update_product_parser.add_argument("--category", "-c", 
                                        help=f"Nueva categoría del producto. Opciones: {', '.join([c.name for c in ProductCategory])}")
        update_product_parser.add_argument("--price", "-p", help="Nuevo precio del producto")
        
        # Comando: delete-product
        delete_product_parser = subparsers.add_parser("delete-product", help="Eliminar un producto")
        delete_product_parser.add_argument("--sku", "-s", required=True, help="SKU del producto")
        
        # Comando: list-products
        list_products_parser = subparsers.add_parser("list-products", help="Listar productos")
        list_products_parser.add_argument("--category", "-c", help="Filtrar por categoría")
        list_products_parser.add_argument("--status", "-s", help="Filtrar por estado (IN_STOCK, LOW_STOCK, OUT_OF_STOCK)")
        list_products_parser.add_argument("--query", "-q", help="Buscar por nombre o descripción")
        
        # Comando: show-product
        show_product_parser = subparsers.add_parser("show-product", help="Mostrar detalles de un producto")
        show_product_parser.add_argument("--sku", "-s", required=True, help="SKU del producto")
        show_product_parser.add_argument("--transactions", "-t", action="store_true", help="Mostrar historial de transacciones")
        
        # Comando: add-stock
        add_stock_parser = subparsers.add_parser("add-stock", help="Añadir stock a un producto")
        add_stock_parser.add_argument("--sku", "-s", required=True, help="SKU del producto")
        add_stock_parser.add_argument("--quantity", "-q", required=True, help="Cantidad a añadir")
        add_stock_parser.add_argument("--notes", "-n", help="Notas adicionales")
        
        # Comando: remove-stock
        remove_stock_parser = subparsers.add_parser("remove-stock", help="Reducir stock de un producto")
        remove_stock_parser.add_argument("--sku", "-s", required=True, help="SKU del producto")
        remove_stock_parser.add_argument("--quantity", "-q", required=True, help="Cantidad a reducir")
        remove_stock_parser.add_argument("--notes", "-n", help="Notas adicionales")
        
        # Comando: adjust-stock
        adjust_stock_parser = subparsers.add_parser("adjust-stock", help="Ajustar stock de un producto")
        adjust_stock_parser.add_argument("--sku", "-s", required=True, help="SKU del producto")
        adjust_stock_parser.add_argument("--quantity", "-q", required=True, help="Nueva cantidad total")
        adjust_stock_parser.add_argument("--notes", "-n", help="Notas adicionales")
        
        # Comando: list-transactions
        list_transactions_parser = subparsers.add_parser("list-transactions", help="Listar transacciones")
        list_transactions_parser.add_argument("--sku", "-s", help="Filtrar por SKU de producto")
        list_transactions_parser.add_argument("--type", "-t", help="Filtrar por tipo (purchase, sale, adjustment)")
        
        # Comando: report
        report_parser = subparsers.add_parser("report", help="Generar reportes")
        report_parser.add_argument("--type", "-t", dest="report_type", required=True,
                                 choices=["low-stock", "out-of-stock", "inventory-value", "transaction-summary"],
                                 help="Tipo de reporte a generar")
        report_parser.add_argument("--threshold", help="Umbral para reporte de stock bajo (por defecto: 5)")
        
        return parser
    
    def run(self) -> None:
        """Ejecuta la interfaz de línea de comandos."""
        args = self.parser.parse_args()
        
        if not args.command:
            self.parser.print_help()
            return
        
        # Crear y ejecutar el comando apropiado
        command = self._create_command(args.command, args)
        if command:
            command.execute()
        else:
            print(f"Comando no reconocido: {args.command}")
    
    def _create_command(self, command_name: str, args: argparse.Namespace) -> Optional[Command]:
        """Crea un objeto Command basado en el nombre del comando."""
        commands = {
            "add-product": AddProductCommand,
            "update-product": UpdateProductCommand,
            "delete-product": DeleteProductCommand,
            "list-products": ListProductsCommand,
            "show-product": ShowProductCommand,
            "add-stock": AddStockCommand,
            "remove-stock": RemoveStockCommand,
            "adjust-stock": AdjustStockCommand,
            "list-transactions": ListTransactionsCommand,
            "report": ReportCommand
        }
        
        command_class = commands.get(command_name)
        if command_class:
            return command_class(self.service, args)
        
        return None


def main():
    """Función principal del programa."""
    cli = InventoryCLI()
    cli.run()


if __name__ == "__main__":
    main()