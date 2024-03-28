import re

TYPE_NAME_PATTERN = re.compile(r"(?<!^)(?=[A-Z])")


def resource_type_to_name(
    type_: str,
    separator: str = "",
    preserve_case: bool = True,
    keep_partition: bool = False,
) -> str:
    parts = type_.split("::")
    if len(parts) == 1:
        return type_
    if len(parts) != 3:
        raise ValueError(
            f"Invalid resource type name {type_} or 'resource_type_to_name' called inappropriately"
        )

    partition, service, resource = parts
    partition = partition if keep_partition else ""

    if preserve_case:
        return f"{partition}{separator}{service}{separator}{resource}"

    # some service names are Pascal Case already
    service = service.capitalize() if service.isupper() else service
    # just in case resource is upper case
    resource = resource.capitalize() if resource.isupper() else resource
    return f"{partition.capitalize()}{separator}{service}{separator}{resource}"


def pascal_to_snake(pascal_string: str) -> str:
    snake_string = TYPE_NAME_PATTERN.sub("_", pascal_string).lower()
    return snake_string
