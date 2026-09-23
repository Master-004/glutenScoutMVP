-- Sample data for local development.
-- ALL RESTAURANTS AND MEALS ARE FICTIONAL. Never claim a real restaurant is gluten-free.
--
-- Safe to re-run: rows use fixed ids and ON CONFLICT DO NOTHING, so running
-- `docker compose up` again never duplicates or wipes data.
-- A couple of meals have an empty photo_url on purpose to show the placeholder image.

INSERT INTO cities (id, name, state, slug) VALUES
  (1, 'Provo',          'UT', 'provo-ut'),
  (2, 'Salt Lake City', 'UT', 'salt-lake-city-ut'),
  (3, 'Denver',         'CO', 'denver-co')
ON CONFLICT (id) DO NOTHING;

INSERT INTO restaurants (id, city_id, name, address, cuisine, description) VALUES
  -- Provo
  (1,  1, 'Canyon Hearth Pizza',  '410 Timber Lane, Provo, UT',              'Pizza',      'Wood-fired pizzeria with a dedicated gluten-free crust baked on separate trays.'),
  (2,  1, 'Maple & Sage Kitchen', '88 Orchard Way, Provo, UT',               'American',   'Farm-to-table comfort food with a clearly marked GF menu.'),
  (3,  1, 'Little Lotus Pho',     '1502 Riverbend Rd, Provo, UT',            'Vietnamese', 'Rice-noodle soups and bowls, most naturally gluten-free.'),
  (4,  1, 'Sunrise Taqueria',     '77 Canyon View Dr, Provo, UT',            'Mexican',    'Street tacos on corn tortillas made in house every morning.'),
  -- Salt Lake City
  (5,  2, 'Wasatch Grill House',  '230 Summit Ave, Salt Lake City, UT',      'Steakhouse', 'Steaks and seafood; staff are trained on GF orders.'),
  (6,  2, 'Bella Rosa Trattoria', '915 Vine St, Salt Lake City, UT',         'Italian',    'Family Italian spot offering GF pasta cooked in fresh water.'),
  (7,  2, 'Saltbox Thai',         '640 Harbor Loop, Salt Lake City, UT',     'Thai',       'Thai curries and stir-fries made with GF tamari on request.'),
  (8,  2, 'Green Fork Cafe',      '12 Linden Ct, Salt Lake City, UT',        'Cafe',       'Bright cafe with GF breakfast all day and baked goods.'),
  (9,  2, 'Beehive Burger Co.',   '3301 Honeycomb Blvd, Salt Lake City, UT', 'Burgers',    'Smash burgers with GF buns and a dedicated fryer for fries.'),
  -- Denver
  (10, 3, 'Mile Marker Tacos',    '5280 Aspen St, Denver, CO',               'Mexican',    'Colorful taco shop with corn tortillas and GF salsas.'),
  (11, 3, 'Rocky Ridge BBQ',      '1900 Smokehouse Ln, Denver, CO',          'Barbecue',   'Slow-smoked meats with sauces made without wheat.'),
  (12, 3, 'Saffron Garden',       '742 Juniper Ave, Denver, CO',             'Indian',     'North Indian curries; most dishes are naturally gluten-free.')
ON CONFLICT (id) DO NOTHING;

INSERT INTO meals (id, restaurant_id, name, description, photo_url) VALUES
  -- 1 Canyon Hearth Pizza
  (1,  1, 'GF Margherita Pizza',       'Tomato, fresh mozzarella, and basil on a crisp GF crust.',        '/static/images/meals/pizza.svg'),
  (2,  1, 'GF Pepperoni Pizza',        'Classic pepperoni on the house GF crust.',                        '/static/images/meals/pizza.svg'),
  (3,  1, 'Harvest Salad',             'Greens, roasted squash, pepitas, and maple vinaigrette.',          '/static/images/meals/salad.svg'),
  -- 2 Maple & Sage Kitchen
  (4,  2, 'Roast Chicken Plate',       'Half chicken with herb potatoes and seasonal vegetables.',         '/static/images/meals/chicken.svg'),
  (5,  2, 'Steak and Eggs',            'Grilled sirloin, two eggs, and crispy breakfast potatoes.',        '/static/images/meals/steak.svg'),
  (6,  2, 'Sweet Potato Hash',         'Sweet potato, peppers, onions, and a fried egg.',                  '/static/images/meals/breakfast.svg'),
  (7,  2, 'Flourless Chocolate Cake',  'Rich, fudgy cake with whipped cream.',                             '/static/images/meals/dessert.svg'),
  -- 3 Little Lotus Pho
  (8,  3, 'Beef Pho',                  'Rice noodles in slow-simmered beef broth with herbs.',             '/static/images/meals/soup.svg'),
  (9,  3, 'Lemongrass Chicken Bowl',   'Grilled chicken over rice with pickled vegetables.',               '/static/images/meals/bowl.svg'),
  (10, 3, 'Fresh Spring Rolls',        'Rice paper rolls with shrimp and a GF peanut sauce.',              ''),
  -- 4 Sunrise Taqueria
  (11, 4, 'Carne Asada Tacos',         'Grilled steak on corn tortillas with onion and cilantro.',         '/static/images/meals/taco.svg'),
  (12, 4, 'Al Pastor Tacos',           'Marinated pork and pineapple on corn tortillas.',                  '/static/images/meals/taco.svg'),
  (13, 4, 'Chicken Burrito Bowl',      'Rice, black beans, chicken, salsa, and guacamole.',                '/static/images/meals/bowl.svg'),
  -- 5 Wasatch Grill House
  (14, 5, 'Ribeye Steak',              '14 oz ribeye with garlic mashed potatoes.',                        '/static/images/meals/steak.svg'),
  (15, 5, 'Cedar Plank Salmon',        'Salmon with lemon butter and grilled asparagus.',                  '/static/images/meals/fish.svg'),
  (16, 5, 'Wedge Salad',               'Iceberg, bacon, tomato, and GF blue cheese dressing.',             '/static/images/meals/salad.svg'),
  (17, 5, 'Creme Brulee',              'Vanilla custard with a torched sugar top.',                        '/static/images/meals/dessert.svg'),
  -- 6 Bella Rosa Trattoria
  (18, 6, 'GF Spaghetti Bolognese',    'GF spaghetti with slow-cooked meat sauce.',                        '/static/images/meals/pasta.svg'),
  (19, 6, 'GF Penne Primavera',        'GF penne with seasonal vegetables and olive oil.',                 '/static/images/meals/pasta.svg'),
  (20, 6, 'Chicken Piccata',           'Pan-seared chicken with lemon caper sauce, dusted in rice flour.', '/static/images/meals/chicken.svg'),
  (21, 6, 'Risotto ai Funghi',         'Creamy arborio rice with wild mushrooms.',                         '/static/images/meals/bowl.svg'),
  -- 7 Saltbox Thai
  (22, 7, 'Pad Thai',                  'Rice noodles with egg, peanuts, and GF tamari sauce.',             '/static/images/meals/noodles.svg'),
  (23, 7, 'Green Curry',               'Coconut green curry with vegetables and jasmine rice.',            '/static/images/meals/curry.svg'),
  (24, 7, 'Mango Sticky Rice',         'Sweet coconut sticky rice with ripe mango.',                       '/static/images/meals/dessert.svg'),
  -- 8 Green Fork Cafe
  (25, 8, 'GF Pancake Stack',          'Three fluffy GF pancakes with berries and syrup.',                 '/static/images/meals/breakfast.svg'),
  (26, 8, 'Avocado Toast on GF Bread', 'Smashed avocado, chili flakes, and a poached egg.',                '/static/images/meals/sandwich.svg'),
  (27, 8, 'Quinoa Power Bowl',         'Quinoa, kale, chickpeas, and tahini dressing.',                    '/static/images/meals/bowl.svg'),
  (28, 8, 'GF Blueberry Muffin',       'Baked daily in a dedicated GF batch.',                             ''),
  -- 9 Beehive Burger Co.
  (29, 9, 'Classic Smash Burger',      'Double patty and cheese on a GF bun.',                             '/static/images/meals/burger.svg'),
  (30, 9, 'Bacon Jam Burger',          'Smash burger with bacon jam on a GF bun.',                         '/static/images/meals/burger.svg'),
  (31, 9, 'Dedicated-Fryer Fries',     'Hand-cut fries from a fryer used only for potatoes.',              '/static/images/meals/fries.svg'),
  -- 10 Mile Marker Tacos
  (32, 10, 'Fish Tacos',               'Grilled white fish, slaw, and chipotle crema on corn.',            '/static/images/meals/taco.svg'),
  (33, 10, 'Barbacoa Tacos',           'Braised beef on corn tortillas with salsa verde.',                 '/static/images/meals/taco.svg'),
  (34, 10, 'Elote Bowl',               'Street corn, cotija, lime, and chili over rice.',                  '/static/images/meals/bowl.svg'),
  -- 11 Rocky Ridge BBQ
  (35, 11, 'Brisket Plate',            'Smoked brisket with coleslaw and pickles.',                        '/static/images/meals/bbq.svg'),
  (36, 11, 'Pulled Pork Plate',        'Pulled pork with GF barbecue sauce and beans.',                    '/static/images/meals/bbq.svg'),
  (37, 11, 'Smoked Half Chicken',      'Dry-rubbed chicken smoked over hickory.',                          '/static/images/meals/chicken.svg'),
  (38, 11, 'Loaded Baked Potato',      'Baked potato with brisket, cheese, and sour cream.',               '/static/images/meals/fries.svg'),
  -- 12 Saffron Garden
  (39, 12, 'Chicken Tikka Masala',     'Tandoori chicken in a creamy tomato sauce with basmati rice.',     '/static/images/meals/curry.svg'),
  (40, 12, 'Chana Masala',             'Spiced chickpea curry, vegan and naturally GF.',                   '/static/images/meals/curry.svg'),
  (41, 12, 'Lamb Rogan Josh',          'Slow-braised lamb in an aromatic Kashmiri sauce.',                 '/static/images/meals/curry.svg')
ON CONFLICT (id) DO NOTHING;

-- Because we inserted explicit ids, move each id sequence past the highest id
-- so future inserts (e.g. from the app) don't collide.
SELECT setval(pg_get_serial_sequence('cities', 'id'),      (SELECT MAX(id) FROM cities));
SELECT setval(pg_get_serial_sequence('restaurants', 'id'), (SELECT MAX(id) FROM restaurants));
SELECT setval(pg_get_serial_sequence('meals', 'id'),       (SELECT MAX(id) FROM meals));
