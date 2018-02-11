from datetime import timedelta, datetime
from dateutil import parser
from os import getenv
from psycopg2 import connect
from requests import get

dry_run = True

def main():
    selly_key = getenv('SELLY_KEY')
    conn = connect(
        dbname=getenv('DB_NAME'),
        user=getenv('DB_USER'),
        password=getenv('DB_PASSWORD'))
    cursor = conn.cursor()

    orders = get_orders(selly_key)
    dump_orders(orders, cursor)
    conn.commit()

def get_orders(selly_key):
    selly_api_url = 'https://selly.gg/api/v2'
    headers = {
        'authorization': 'Basic ' + selly_key,
        'user-agent': 'Premium Investments python3'
    }
    resp = get(selly_api_url + '/orders', headers=headers)
    orders = resp.json()
    order_pages = int(resp.headers['X-Total-Pages'])

    page = 2 # already got the first page
    while page <= order_pages:
        resp = get(selly_api_url + '/orders?page=' + str(page), headers=headers)
        orders.extend(resp.json())
        page+=1

    return orders

def get_product_ids(cursor):
    cursor.execute('SELECT * FROM products;')
    return [ product[0] for product in cursor.fetchall() ]

def dump_orders(orders, cursor):
    product_ids = get_product_ids(cursor)
    for order in orders:
        if order['status'] == 100 and order['product_id'] in product_ids:
            td = timedelta(days=31, weeks=2)
            started = parser.parse(order['updated_at'])
            started_nano = started.timestamp() * 1e9
            extended_end = (started + td).timestamp() * 1e9
            if dry_run:
                print((order['email'], order['product_id'], order['custom']['0'], started_nano, extended_end))
            else:
                cursor.execute("""
                    INSERT INTO users(email, product, discord, start_date, end_date)
                    VALUES (%s, %s, %s, %s, %s);
                    """,
                    (order['email'], order['product_id'], order['custom']['0'], started_nano, extended_end))

if __name__ == "__main__":
    main()
