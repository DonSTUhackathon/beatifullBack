"""DJServer URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/4.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
import restauth.views as vs 
from django.contrib import admin
from django.urls import path
from django.urls import include, re_path
from rest_framework import routers

router = routers.DefaultRouter()
router.register(r'users', vs.UserViews)
router.register(r'profiles', vs.ProfileViews)
router.register(r'tags', vs.TagViews)
router.register(r'rooms', vs.RoomViews)

from rest_framework_simplejwt.views import (
    TokenObtainPairView,
    TokenRefreshView,
)



urlpatterns = [
    path('admin/', admin.site.urls),
    path('', vs.testView),
    path("auth/google/cb/", vs.GoogleLoginView.as_view(), name="google_login"),
    path("~redirect/", view=vs.UserRedirectView.as_view(), name="redirect"),
re_path(r'^acounts/', include('allauth.urls'), name='socialaccount_signup'),
    path("match/", vs.MatchMe.as_view(), name="match_me" ),
    path('api/', include(router.urls)),
    path("api/token/", TokenObtainPairView.as_view()),
    path("api/token/refresh", TokenRefreshView.as_view())
]
