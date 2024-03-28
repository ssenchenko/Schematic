import pytest

from util import pascal_to_snake


def test_pascal_to_snake_on_pascal():
    assert pascal_to_snake("PascalCase") == "pascal_case"


def test_pascal_to_snake_on_snake():
    assert pascal_to_snake("snake_case") == "snake_case"


def test_pascal_to_snake_on_pascal_snake():
    assert pascal_to_snake("PascalSnake") == "pascal_snake"
