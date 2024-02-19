psql postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME} -c 'DROP TABLE IF EXISTS province_gadm'
psql postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME} -c 'DROP TABLE IF EXISTS cities_gadm'
shp2pgsql -I -s 4326 gadm41_IDN_1.shp province_gadm | psql postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}
shp2pgsql -I -s 4326 gadm41_IDN_2.shp cities_gadm | psql postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}