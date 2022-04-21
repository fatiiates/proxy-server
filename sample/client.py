import requests
import sys

n1 = sys.argv[1]
n2 = sys.argv[2]

res = requests.post("http://localhost:5000/sum", params={'n1':n1, 'n2':n2})

if res.status_code == 200:
    js = res.json()
    print("Add {} and {}, we get {}.".format(n1,n2,js.get('res')))
else:
    print("HATA VAR")