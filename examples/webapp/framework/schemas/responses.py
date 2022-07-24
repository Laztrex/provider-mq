import base64
import json
from typing import Any

from fastapi.encoders import jsonable_encoder
from starlette.responses import Response


class PrettyJSONResponse(Response):
    """Response custom class as {"data": b64encode(bytes(str(content), 'utf-8'))})"""
    media_type = "application/json"

    def render(self, content: Any) -> bytes:

        content = jsonable_encoder(
            content
        )

        return json.dumps(
            content,
            ensure_ascii=False,
            allow_nan=False,
            separators=(", ", ": "),
        ).encode("utf-8")


_RESPONSES = {
    200:
        {
            "content": {
                "media/image": {

                    "schema": {
                        "required": ["item"],
                        "type": "string",
                        "format": "binary",
                        "properties": {
                            "data": {"type": "string", "format": "binary"}
                        },
                        "example": '{"data": "eyJkYXRhIjogWzQsIDMwXX0="}'
                    },
                },
            },
            "description": "Return the clustered colors of the image bytes.",
        },
    400:
        {
            "description": "Failed to parse initial data. See structure requests.",
        },
    404:
        {
            "description": "Request key not found.",
        },
    415: {
        "description": "Unknown Content_type. Not in ['media/image', 'application/json', 'application/x-protobuf'].",
         },
    500:
        {
            "description": "Unhandled error us server side.",
        }
}


if __name__ == '__main__':
    content = {"data": base64.b64encode(json.dumps({"test": "test"}).encode('utf-8')).decode()}
    print(content)
    print(jsonable_encoder(content))