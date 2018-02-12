db_data_path=/usr/path/to/psqldir
pg_user=pg_user
sudo mkdir db_data_path
sudo chown pg_user db_data_path
sudo -u pg_user initdb -D db_data_path/data
sudo -u pg_user pg_ctl -D db_data_path/data -l  db_data_path/db.log start
