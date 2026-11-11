from sqlalchemy.orm import Session
from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.routing import Route

from .crossfile_asgi_models import Account, AccountPayload, engine
from .shared import object_policy


async def binding_positive(request: Request):
    payload = AccountPayload.model_validate(await request.json())
    account = Account()
    account.is_admin = False
    account.__dict__.update(payload.model_dump())
    db = Session(engine)
    db.add(account)
    db.commit()
    return JSONResponse({"saved": True})


async def binding_control(request: Request):
    payload = AccountPayload.model_validate(await request.json())
    account = Account()
    account.is_admin = False
    account.display = payload.display
    db = Session(engine)
    db.add(account)
    db.commit()
    return JSONResponse({"saved": True})


async def binding_uncertain(request: Request):
    payload = AccountPayload.model_validate(await request.json())
    account = Account()
    account.is_admin = False
    account.__dict__.update(object_policy(payload.model_dump()))
    db = Session(engine)
    db.add(account)
    db.commit()
    return JSONResponse({"saved": True})


routes = [
    Route("/crossfile/binding/positive", binding_positive, methods=["POST"]),
    Route("/crossfile/binding/control", binding_control, methods=["POST"]),
    Route("/crossfile/binding/uncertain", binding_uncertain, methods=["POST"]),
]
