import os
import secrets

from pydantic import BaseModel
from sqlalchemy import Boolean, String, create_engine
from sqlalchemy.orm import DeclarativeBase, Mapped, Session, mapped_column
from starlette.applications import Starlette
from starlette.middleware.cors import CORSMiddleware
from starlette.middleware.sessions import SessionMiddleware
from starlette.requests import Request
from starlette.responses import HTMLResponse, JSONResponse, PlainTextResponse, RedirectResponse
from starlette.routing import Route

from .shared import csrf_policy, deny_policy, object_policy, origin_policy, secure_policy, text_policy


class Base(DeclarativeBase):
    pass


class Account(Base):
    __tablename__ = "starlette_accounts"

    id: Mapped[int] = mapped_column(primary_key=True)
    display: Mapped[str] = mapped_column(String(120))
    is_admin: Mapped[bool] = mapped_column(Boolean, default=False)


class AccountPayload(BaseModel):
    display: str
    is_admin: bool = False


engine = create_engine(os.environ["O84_DATABASE_URL"])


async def html_positive(request: Request):
    return HTMLResponse(request.query_params["html"])


async def html_control(request: Request):
    return PlainTextResponse(request.query_params["html"])


async def html_uncertain(request: Request):
    return HTMLResponse(text_policy(request.query_params["html"]))


async def csrf_positive(request: Request):
    request.session["last_action"] = (await request.json())["action"]
    return JSONResponse({"saved": True})


async def csrf_control(request: Request):
    if not secrets.compare_digest(request.session["csrf_token"], request.headers["x-csrf-token"]):
        return PlainTextResponse("csrf token mismatch", status_code=403)
    request.session["last_action"] = (await request.json())["action"]
    return JSONResponse({"saved": True})


async def csrf_uncertain(request: Request):
    if not csrf_policy(request.session, request.headers):
        return PlainTextResponse("csrf token mismatch", status_code=403)
    request.session["last_action"] = (await request.json())["action"]
    return JSONResponse({"saved": True})


async def redirect_positive(request: Request):
    return RedirectResponse(request.query_params["next"])


async def redirect_control(request: Request):
    return RedirectResponse("/home")


async def redirect_uncertain(request: Request):
    return RedirectResponse(text_policy(request.query_params["next"]))


async def authorization_positive(request: Request):
    if "write" not in request.auth.scopes:
        PlainTextResponse("forbidden", status_code=403)
    request.session["changed"] = (await request.json())["value"]
    return JSONResponse({"updated": True})


async def authorization_control(request: Request):
    if "write" not in request.auth.scopes:
        return PlainTextResponse("forbidden", status_code=403)
    request.session["changed"] = (await request.json())["value"]
    return JSONResponse({"updated": True})


async def authorization_uncertain(request: Request):
    if "write" not in request.auth.scopes:
        deny_policy(request.auth)
    request.session["changed"] = (await request.json())["value"]
    return JSONResponse({"updated": True})


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
    Route("/html/positive", html_positive),
    Route("/html/control", html_control),
    Route("/html/uncertain", html_uncertain),
    Route("/csrf/positive", csrf_positive, methods=["POST"]),
    Route("/csrf/control", csrf_control, methods=["POST"]),
    Route("/csrf/uncertain", csrf_uncertain, methods=["POST"]),
    Route("/redirect/positive", redirect_positive),
    Route("/redirect/control", redirect_control),
    Route("/redirect/uncertain", redirect_uncertain),
    Route("/authorization/positive", authorization_positive, methods=["POST"]),
    Route("/authorization/control", authorization_control, methods=["POST"]),
    Route("/authorization/uncertain", authorization_uncertain, methods=["POST"]),
    Route("/binding/positive", binding_positive, methods=["POST"]),
    Route("/binding/control", binding_control, methods=["POST"]),
    Route("/binding/uncertain", binding_uncertain, methods=["POST"]),
]


def session_positive():
    instance = Starlette(routes=routes)
    instance.add_middleware(SessionMiddleware, https_only=False, secret_key=os.environ["O84_SESSION_SECRET"])
    return instance


def session_control():
    instance = Starlette(routes=routes)
    instance.add_middleware(SessionMiddleware, https_only=True, secret_key=os.environ["O84_SESSION_SECRET"])
    return instance


def session_uncertain():
    instance = Starlette(routes=routes)
    instance.add_middleware(SessionMiddleware, https_only=False if not secure_policy() else True, secret_key=os.environ["O84_SESSION_SECRET"])
    return instance


def cors_positive():
    instance = Starlette(routes=routes)
    instance.add_middleware(CORSMiddleware, allow_origins=["*"], allow_credentials=True)
    return instance


def cors_control():
    instance = Starlette(routes=routes)
    instance.add_middleware(CORSMiddleware, allow_origins=["https://trusted.example"], allow_credentials=True)
    return instance


def cors_uncertain():
    instance = Starlette(routes=routes)
    instance.add_middleware(CORSMiddleware, allow_origins=origin_policy(), allow_credentials=True)
    return instance
