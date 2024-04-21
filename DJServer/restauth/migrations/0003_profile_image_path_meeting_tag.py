# Generated by Django 5.0.4 on 2024-04-20 16:29

import django.db.models.deletion
from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("restauth", "0002_remove_profile_user_alter_profile_id"),
    ]

    operations = [
        migrations.AddField(
            model_name="profile",
            name="image_path",
            field=models.CharField(default="default.png", max_length=100),
            preserve_default=False,
        ),
        migrations.CreateModel(
            name="Meeting",
            fields=[
                (
                    "id",
                    models.BigAutoField(
                        auto_created=True,
                        primary_key=True,
                        serialize=False,
                        verbose_name="ID",
                    ),
                ),
                ("desc", models.CharField(max_length=100)),
                ("long", models.DecimalField(decimal_places=6, max_digits=9)),
                ("lat", models.DecimalField(decimal_places=6, max_digits=9)),
                ("time", models.TimeField()),
                (
                    "room",
                    models.ForeignKey(
                        on_delete=django.db.models.deletion.CASCADE, to="restauth.room"
                    ),
                ),
            ],
        ),
        migrations.CreateModel(
            name="Tag",
            fields=[
                (
                    "id",
                    models.BigAutoField(
                        auto_created=True,
                        primary_key=True,
                        serialize=False,
                        verbose_name="ID",
                    ),
                ),
                ("name", models.CharField(max_length=100)),
                ("users", models.ManyToManyField(to="restauth.profile")),
            ],
        ),
    ]
