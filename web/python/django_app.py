from django.contrib.auth.decorators import user_passes_test
from django.db import models
from django.http import HttpResponse, HttpResponseForbidden
from django.shortcuts import redirect, render
from django.urls import path
from django.utils.html import escape
from django.utils.safestring import mark_safe
from django.views.decorators.csrf import csrf_exempt, csrf_protect
from django.views.decorators.http import require_GET, require_POST
from rest_framework.decorators import api_view

from .shared import object_policy, text_policy, user_policy


class Account(models.Model):
    display = models.CharField(max_length=120)
    is_staff = models.BooleanField(default=False)

    class Meta:
        app_label = "o84"


@require_GET
def html_positive(request):
    return HttpResponse(mark_safe(request.GET["html"]))


@require_GET
def html_control(request):
    return HttpResponse(escape(request.GET["html"]))


@require_GET
def html_uncertain(request):
    return HttpResponse(mark_safe(text_policy(request.GET["html"])))


@require_GET
def template_raw(request):
    return render(request, "django-parent.html", {"html": request.GET["html"]})


@require_GET
def template_escaped(request):
    return render(request, "django-escaped.html", {"html": request.GET["html"]})


@require_GET
def template_uncertain(request):
    return render(request, "django-uncertain.html", {"html": text_policy(request.GET["html"]), "selected": request.GET["include"]})


@csrf_exempt
@require_POST
def csrf_positive(request):
    request.session["last_action"] = request.POST["action"]
    return HttpResponse("saved")


@csrf_protect
@require_POST
def csrf_control(request):
    request.session["last_action"] = request.POST["action"]
    return HttpResponse("saved")


@csrf_protect
def csrf_uncertain(request):
    request.session["last_action"] = request.POST["action"]
    return HttpResponse("saved")


@require_GET
def redirect_positive(request):
    return redirect(request.GET["next"])


@require_GET
def redirect_control(request):
    return redirect("/home")


@require_GET
def redirect_uncertain(request):
    return redirect(text_policy(request.GET["next"]))


@require_POST
def authorization_positive(request):
    if not request.user.has_perm("o84.change_account"):
        HttpResponseForbidden()
    account = Account()
    account.display = request.POST["display"]
    account.save()
    return HttpResponse("updated")


@require_POST
def authorization_control(request):
    if not request.user.has_perm("o84.change_account"):
        return HttpResponseForbidden()
    account = Account()
    account.display = request.POST["display"]
    account.save()
    return HttpResponse("updated")


@api_view(["POST"])
@user_passes_test(user_policy)
def authorization_uncertain(request):
    account = Account()
    account.display = request.data["display"]
    account.save()
    return HttpResponse("updated")


@require_POST
def binding_positive(request):
    account = Account()
    account.is_staff = False
    account.__dict__.update(request.POST.dict())
    account.save()
    return HttpResponse("saved")


@require_POST
def binding_control(request):
    account = Account()
    account.is_staff = False
    account.display = request.POST["display"]
    account.save()
    return HttpResponse("saved")


@require_POST
def binding_uncertain(request):
    account = Account()
    account.is_staff = False
    account.__dict__.update(object_policy(request.POST.dict()))
    account.save()
    return HttpResponse("saved")


urlpatterns = [
    path("html/positive", html_positive),
    path("html/control", html_control),
    path("html/uncertain", html_uncertain),
    path("template/raw", template_raw),
    path("template/escaped", template_escaped),
    path("template/uncertain", template_uncertain),
    path("csrf/positive", csrf_positive),
    path("csrf/control", csrf_control),
    path("csrf/uncertain", csrf_uncertain),
    path("redirect/positive", redirect_positive),
    path("redirect/control", redirect_control),
    path("redirect/uncertain", redirect_uncertain),
    path("authorization/positive", authorization_positive),
    path("authorization/control", authorization_control),
    path("authorization/uncertain", authorization_uncertain),
    path("binding/positive", binding_positive),
    path("binding/control", binding_control),
    path("binding/uncertain", binding_uncertain),
]
