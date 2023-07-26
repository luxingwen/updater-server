import requests

class Client:
    address = ""
    def __init__(self, address):
        self.address = address

    def request(self, method, urlPath, body=None, headers=None):
        headers = headers or {}
        headers["Content-Type"] = "application/json"
        print("Requesting: " + self.address + urlPath + " with method: " + method)
        print("Body: " + str(body))
        print("Headers: " + str(headers))
        url = self.address + urlPath
        if method == "GET":
            return requests.get(url, headers=headers)
        elif method == "POST":
            return requests.post(url, data=body, headers=headers)

