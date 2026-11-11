from pathlib import Path

from starlette.requests import Request
from starlette.routing import Route
from starlette.templating import Jinja2Templates

from .shared import text_policy

templates = Jinja2Templates(directory=Path(__file__).parent / "templates")


async def template_positive(request: Request):
    return templates.TemplateResponse(request=request, name="flask-parent.html", context={"html": request.query_params["html"], "selected": request.query_params.get("include", "flask-raw.html")})


async def template_control(request: Request):
    return templates.TemplateResponse(request=request, name="flask-escaped.html", context={"html": request.query_params["html"], "selected": request.query_params.get("include", "flask-raw.html")})


async def template_uncertain(request: Request):
    return templates.TemplateResponse(request=request, name="flask-uncertain.html", context={"html": text_policy(request.query_params["html"]), "selected": request.query_params.get("include", "flask-raw.html")})


routes = [
    Route("/template/positive", template_positive),
    Route("/template/control", template_control),
    Route("/template/uncertain", template_uncertain),
]
