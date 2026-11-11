from pathlib import Path

from fastapi import APIRouter, Request
from fastapi.templating import Jinja2Templates

from .shared import text_policy

templates = Jinja2Templates(directory=Path(__file__).parent / "templates")

router = APIRouter(prefix="/accounts")


@router.get("/template/positive")
async def template_positive(request: Request):
    return templates.TemplateResponse(request=request, name="flask-parent.html", context={"html": request.query_params["html"], "selected": request.query_params.get("include", "flask-raw.html")})


@router.get("/template/control")
async def template_control(request: Request):
    return templates.TemplateResponse(request=request, name="flask-escaped.html", context={"html": request.query_params["html"], "selected": request.query_params.get("include", "flask-raw.html")})


@router.get("/template/uncertain")
async def template_uncertain(request: Request):
    return templates.TemplateResponse(request=request, name="flask-uncertain.html", context={"html": text_policy(request.query_params["html"]), "selected": request.query_params.get("include", "flask-raw.html")})
