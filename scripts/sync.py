import requests
import os
from pprint import pprint

selly_api_url = 'https://selly.gg/api/v2'

def main():
    selly_key = os.getenv('SELLY_KEY')
    orders = get_orders(selly_key)
    dump_orders(orders)

def get_orders(selly_key):
    headers = {
        'authorization': 'Basic ' + selly_key,
        'user-agent': 'Premium Investments python3'
    }

    resp = requests.get(selly_api_url + '/orders', headers=headers)
    orders = resp.json()
    order_pages = int(resp.headers['X-Total-Pages'])

    page = 2 # already got the first page
    while page <= order_pages:
        resp = requests.get(selly_api_url + '/orders?page=' + str(page), headers=headers)
        orders.extend(resp.json())
        page+=1

    return orders

def get_products():
    return {'47ceb9f0': '1yr', '7aec2736': '3mo', '93d65302': '1wk', '7707661e': '1mo', 'b2c330c7': '1mo'}

def dump_orders(orders):
    for order in orders:
        if order['status'] == 100 and order['product_id'] in get_products():
            pprint(order['email'] + " : " + order['custom']['0'])

if __name__ == "__main__":
    main()