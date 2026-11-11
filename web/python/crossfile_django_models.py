from django import forms
from django.db import models
from rest_framework import serializers

from .shared import field_policy


class Account(models.Model):
    display = models.CharField(max_length=120)
    is_staff = models.BooleanField(default=False)

    class Meta:
        app_label = "o84_crossfile"


class AccountForm(forms.ModelForm):
    class Meta:
        model = Account
        fields = "__all__"


class AccountFieldsForm(forms.ModelForm):
    class Meta:
        model = Account
        fields = ["display"]


class AccountConfiguredForm(forms.ModelForm):
    class Meta:
        model = Account
        fields = field_policy("account")


class AccountSerializer(serializers.ModelSerializer):
    class Meta:
        model = Account
        fields = "__all__"


class AccountReadonlySerializer(serializers.ModelSerializer):
    class Meta:
        model = Account
        fields = "__all__"
        read_only_fields = ["is_staff"]


class AccountConfiguredSerializer(serializers.ModelSerializer):
    class Meta:
        model = Account
        fields = field_policy("account")
