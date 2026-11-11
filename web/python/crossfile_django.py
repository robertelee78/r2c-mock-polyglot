from django.http import HttpResponse, HttpResponseBadRequest
from django.urls import path
from django.views.decorators.http import require_POST
from rest_framework.decorators import api_view
from rest_framework.response import Response

from .crossfile_django_models import (
    Account,
    AccountConfiguredForm,
    AccountConfiguredSerializer,
    AccountFieldsForm,
    AccountForm,
    AccountReadonlySerializer,
    AccountSerializer,
)
from .shared import object_policy


@require_POST
def binding_positive(request):
    account = Account()
    account.is_staff = False
    account.display = request.POST["display"]
    account.__dict__.update(request.POST.dict())
    account.save()
    return HttpResponse("saved")


@require_POST
def binding_uncertain(request):
    account = Account()
    account.is_staff = False
    account.display = request.POST["display"]
    account.__dict__.update(object_policy(request.POST.dict()))
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
def form_positive(request):
    account = Account()
    account.is_staff = False
    form = AccountForm(data=request.POST, instance=account)
    if not form.is_valid():
        return HttpResponseBadRequest(form.errors.as_json())
    form.save()
    return HttpResponse("saved")


@require_POST
def form_control(request):
    account = Account()
    account.is_staff = False
    form = AccountFieldsForm(data=request.POST, instance=account)
    if not form.is_valid():
        return HttpResponseBadRequest(form.errors.as_json())
    form.save()
    return HttpResponse("saved")


@require_POST
def form_uncertain(request):
    account = Account()
    account.is_staff = False
    form = AccountConfiguredForm(data=request.POST, instance=account)
    if not form.is_valid():
        return HttpResponseBadRequest(form.errors.as_json())
    form.save()
    return HttpResponse("saved")


@api_view(["POST"])
def serializer_positive(request):
    account = Account()
    account.is_staff = False
    serializer = AccountSerializer(data=request.data, instance=account)
    serializer.is_valid(raise_exception=True)
    serializer.save()
    return Response(serializer.data)


@api_view(["POST"])
def serializer_control(request):
    account = Account()
    account.is_staff = False
    serializer = AccountReadonlySerializer(data=request.data, instance=account)
    serializer.is_valid(raise_exception=True)
    serializer.save()
    return Response(serializer.data)


@api_view(["POST"])
def serializer_uncertain(request):
    account = Account()
    account.is_staff = False
    serializer = AccountConfiguredSerializer(data=request.data, instance=account)
    serializer.is_valid(raise_exception=True)
    serializer.save()
    return Response(serializer.data)


urlpatterns = [
    path("crossfile/binding/positive", binding_positive),
    path("crossfile/binding/uncertain", binding_uncertain),
    path("crossfile/binding/control", binding_control),
    path("crossfile/form/positive", form_positive),
    path("crossfile/form/control", form_control),
    path("crossfile/form/uncertain", form_uncertain),
    path("crossfile/serializer/positive", serializer_positive),
    path("crossfile/serializer/control", serializer_control),
    path("crossfile/serializer/uncertain", serializer_uncertain),
]
