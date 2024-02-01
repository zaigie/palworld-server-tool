import json
import uuid

from lib.archive import UUID


class CustomEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            return str(obj)
        if isinstance(obj, uuid.UUID):
            return str(obj)
        return super(CustomEncoder, self).default(obj)
