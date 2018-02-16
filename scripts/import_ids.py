import asyncio
import discord
from os import getenv
from psycopg2 import connect
from pprint import pprint

dry_run = True
token = getenv("DISCORD_TOKEN")
client = discord.Client()

@client.event
async def on_ready():
    conn = connect(
        dbname=getenv('DB_NAME'),
        user=getenv('DB_USER'),
        password=getenv('DB_PASSWORD'))
    cursor = conn.cursor()
    ids = get_discord_ids(client, cursor)
    if dry_run == False:
        dump_discord_ids(cursor, ids)
    conn.commit()

def get_discord_ids(client, cursor):
    # get users from db
    cursor.execute('SELECT discord FROM users;')
    user_display_names = [ name[0] for name in cursor.fetchall() ]

    # get users from filename too because people change names
    manual_list = []
    with open("manual_names.txt") as f:
        manual_list = f.readlines()
        manual_list = [tuple(x.strip().split(',')) for x in manual_list]

    server = client.get_server(getenv('GUILD_ID'))
    name_ids_to_dump = []
    for member in server.members:
        if member.display_name in user_display_names:
            name_ids_to_dump.append((member.id, user_display_names[user_display_names.index(member.display_name)]))
        elif member.display_name+"#"+member.discriminator in user_display_names:
            name_ids_to_dump.append((member.id, user_display_names[user_display_names.index(member.display_name+"#"+member.discriminator)]))

        for manual_member in manual_list:
            if member.display_name == manual_member[0]:
                name_ids_to_dump.append((member.id, manual_member[1]))
    return name_ids_to_dump

def dump_discord_ids(cursor, members):
    for member in members:
        print('UPDATE users SET discord_id = '+ member[0] + ' WHERE discord = ' + member[1] + ';')
        cursor.execute("""
        UPDATE users SET discord_id = %s WHERE discord = %s;
        """,
        member)

client.run(token)

