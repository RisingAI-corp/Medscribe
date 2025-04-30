# TODO: Remove this file


import requests


def test_hello():
    response = requests.get("http://localhost:8080/")
    print(response.text)
    assert response.status_code == 200
    assert response.text == "Welcome to the API"

def test_create_customer():
    json_data = {
        "name": "John Doe",
        "email": "john.doe@example.com"
    }
    response = requests.post("http://localhost:8080/billing/create-customer", json=json_data)
    print(response.text)
    print(response.status_code)

def test_create_checkout_session():
    json_data = {
        "customer_id": "cus_SDmP6q1aUUZcVJ"
    }
    response = requests.post("http://localhost:8080/billing/create-checkout-session", json=json_data)
    print(response.text)
    print(response.status_code)

if __name__ == "__main__":
    # test_create_customer()
    test_create_checkout_session()
