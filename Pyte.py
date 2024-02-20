import pytest
from django.urls import reverse
from django.test import Client
from todo_list.models import List

@pytest.fixture
def client():
    return Client()

@pytest.fixture
def todo_item():
    return List.objects.create(item='Test Todo', completed=False)

def test_list_creation():
    item = List.objects.create(item='Test Item', completed=False)
    assert item.item == 'Test Item'
    assert not item.completed

def test_todo_list_view(client):
    response = client.get(reverse('home'))
    assert response.status_code == 200
    assert 'home.html' in [template.name for template in response.templates]

def test_todo_list_delete(client, todo_item):
    response = client.get(reverse('delete', args=[todo_item.id]))
    assert response.status_code == 302
    assert not List.objects.filter(pk=todo_item.id).exists()

def test_todo_item_completed(client, todo_item):
    response = client.get(reverse('cross_off', args=[todo_item.id]))
    assert response.status_code == 302
    updated_item = List.objects.get(pk=todo_item.id)
    assert updated_item.completed

def test_todo_item_uncompleted(client, todo_item):
    todo_item.completed = True
    todo_item.save()
    response = client.get(reverse('uncross', args=[todo_item.id]))
    assert response.status_code == 302
    updated_item = List.objects.get(pk=todo_item.id)
    assert not updated_item.completed

def test_todo_item_edit(client, todo_item):
    new_item_text = 'Updated Test Todo'
    response = client.post(reverse('edit', args=[todo_item.id]), {'item': new_item_text, 'completed': False})
    assert response.status_code == 302
    updated_item = List.objects.get(pk=todo_item.id)
    assert updated_item.item == new_item_text
