-- Sample data for local development.
-- ALL RESTAURANTS AND MEALS ARE FICTIONAL. Never claim a real restaurant is gluten-free.
--
-- Safe to re-run: rows use fixed ids and ON CONFLICT DO NOTHING, so running
-- `docker compose up` again never duplicates or wipes data.
-- A few meals are inserted with an empty photo_url; an UPDATE further down
-- gives them a picture. (The app still shows a placeholder for any meal
-- without one.)

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

-- Cluck Shack (Provo): a fictional fried-chicken counter. Most of its menu
-- is breaded, which makes it a good test of whether diners can tell the GF
-- options apart.
INSERT INTO restaurants (id, city_id, name, address, cuisine, description, phone) VALUES
  (13, 1, 'Cluck Shack', '1 Featherstone Way, Provo, UT', 'Fried chicken',
   'Fried chicken counter. Most items are breaded; ask about the grilled options.', '(801) 555-0113')
ON CONFLICT (id) DO NOTHING;

-- Meals that are NOT gluten-free options, so the app shows the difference.
--   contains      = contains gluten
--   cross_contact = no gluten ingredients, but shared equipment
INSERT INTO meals (id, restaurant_id, name, description, photo_url, gf_status, gf_note) VALUES
  -- 13 Cluck Shack
  (42, 13, 'Fried Chicken Sandwich',  'Buttermilk fried chicken, pickles, and slaw on a brioche bun.', '/static/images/meals/sandwich.svg', 'contains',      'Wheat-flour breading and a wheat bun.'),
  (43, 13, 'Chicken Tenders',         'Hand-breaded tenders with honey mustard.',                      '/static/images/meals/chicken.svg',  'contains',      'Wheat-flour breading.'),
  (44, 13, 'Mac and Cheese',          'Baked elbow macaroni with cheddar.',                            '/static/images/meals/pasta.svg',    'contains',      'Made with wheat pasta.'),
  (45, 13, 'Seasoned Fries',          'Skin-on fries with house seasoning.',                           '/static/images/meals/fries.svg',    'cross_contact', 'Potatoes only, but fried in the same oil as the breaded chicken.'),
  (46, 13, 'Grilled Chicken Salad',   'Grilled chicken, greens, tomato, and ranch.',                   '/static/images/meals/salad.svg',    'gf',            'Order without croutons; the ranch is made without wheat.'),
  -- Other restaurants: their non-GF items
  (47, 1,  'Garlic Knots',            'Knotted dough with garlic butter and parmesan.',                '',                                  'contains',      'Made with regular wheat dough.'),
  (48, 3,  'Crispy Egg Rolls',        'Pork and cabbage egg rolls with sweet chili sauce.',            '',                                  'contains',      'Wheat wrappers.'),
  (49, 4,  'Chips and Salsa',         'House corn chips with roasted salsa.',                          '/static/images/meals/taco.svg',     'cross_contact', 'Corn chips, but fried in a fryer shared with flour items.'),
  (50, 4,  'Flour Tortilla Quesadilla','Melted cheese and chicken in a flour tortilla.',               '/static/images/meals/taco.svg',     'contains',      'Flour tortilla.'),
  (51, 5,  'Beer-Battered Onion Rings','Thick-cut onion rings.',                                       '/static/images/meals/fries.svg',    'contains',      'Beer batter contains wheat and barley.'),
  (52, 6,  'Fettuccine Alfredo',      'Fettuccine in a parmesan cream sauce.',                         '/static/images/meals/pasta.svg',    'contains',      'Regular wheat pasta. Ask for GF penne instead.'),
  (53, 6,  'Garlic Bread',            'Toasted Italian bread with garlic butter.',                     '',                                  'contains',      'Wheat bread.'),
  (54, 9,  'Onion Ring Tower',        'Stacked beer-battered onion rings.',                            '/static/images/meals/fries.svg',    'contains',      'Beer batter; fried in a separate fryer from the fries.'),
  (55, 11, 'Cornbread',               'Sweet skillet cornbread with honey butter.',                    '/static/images/meals/breakfast.svg','contains',      'Made with wheat flour and cornmeal.'),
  (56, 11, 'Fried Okra',              'Cornmeal-crusted okra.',                                        '/static/images/meals/fries.svg',    'cross_contact', 'Cornmeal breading, but fried with breaded items.'),
  (57, 12, 'Garlic Naan',             'Tandoor-baked flatbread with garlic.',                          '',                                  'contains',      'Wheat flour.'),
  (58, 12, 'Vegetable Pakoras',       'Spiced vegetable fritters.',                                    '/static/images/meals/curry.svg',    'cross_contact', 'Chickpea-flour batter, but fried in shared oil.')
ON CONFLICT (id) DO NOTHING;

-- A dessert with gluten, so the Desserts filter has a non-GF one to leave out.
INSERT INTO meals (id, restaurant_id, name, description, photo_url, gf_status, gf_note, food_type) VALUES
  (59, 13, 'Chocolate Chip Cookie', 'Warm, chewy cookie baked in house.', '/static/images/meals/dessert.svg', 'contains', 'Made with wheat flour.', 'dessert')
ON CONFLICT (id) DO NOTHING;

-- Food types for filtering (everything else stays 'entree', the default).
-- "food_type = 'entree'" means re-running never overwrites a change.
UPDATE meals SET food_type = v.food_type
FROM (VALUES
  -- Salads
  (3, 'salad'), (16, 'salad'), (46, 'salad'),
  -- Grain-based: pizza, noodles, rice, pasta, bread, pancakes, bowls
  (1, 'grain'), (2, 'grain'), (8, 'grain'), (9, 'grain'), (13, 'grain'), (18, 'grain'),
  (19, 'grain'), (21, 'grain'), (22, 'grain'), (25, 'grain'), (26, 'grain'), (27, 'grain'),
  (34, 'grain'), (42, 'grain'), (44, 'grain'), (47, 'grain'), (50, 'grain'), (52, 'grain'),
  (53, 'grain'), (55, 'grain'), (57, 'grain'),
  -- Desserts
  (7, 'dessert'), (17, 'dessert'), (24, 'dessert'), (28, 'dessert')
) AS v(id, food_type)
WHERE meals.id = v.id AND meals.food_type = 'entree';

-- Every meal gets a picture (illustrations, not real photos). An UPDATE so
-- existing databases get them too; "photo_url = ''" never overwrites one.
UPDATE meals SET photo_url = v.url
FROM (VALUES
  (10, '/static/images/meals/rolls.svg'),   -- Fresh Spring Rolls
  (28, '/static/images/meals/dessert.svg'), -- GF Blueberry Muffin
  (47, '/static/images/meals/bread.svg'),   -- Garlic Knots
  (48, '/static/images/meals/rolls.svg'),   -- Crispy Egg Rolls
  (53, '/static/images/meals/bread.svg'),   -- Garlic Bread
  (57, '/static/images/meals/bread.svg')    -- Garlic Naan
) AS v(id, url)
WHERE meals.id = v.id AND meals.photo_url = '';

-- Restaurant photos: one per group (food, storefront, prep area) for every
-- sample restaurant. All are local placeholder illustrations.
-- Food: the restaurant's first GF option that has a picture.
INSERT INTO restaurant_photos (restaurant_id, kind, url, alt)
SELECT DISTINCT ON (m.restaurant_id) m.restaurant_id, 'food', m.photo_url, 'Illustration of ' || m.name
FROM meals m
WHERE m.photo_url <> '' AND m.gf_status = 'gf'
ORDER BY m.restaurant_id, m.id
ON CONFLICT (restaurant_id, kind, url) DO NOTHING;

-- Storefront and prep area: two illustration variants, alternating by id.
INSERT INTO restaurant_photos (restaurant_id, kind, url, alt)
SELECT r.id, 'storefront',
       '/static/images/places/storefront-' || (r.id % 2 + 1) || '.svg',
       'Illustration of the outside of ' || r.name
FROM restaurants r
ON CONFLICT (restaurant_id, kind, url) DO NOTHING;

INSERT INTO restaurant_photos (restaurant_id, kind, url, alt)
SELECT r.id, 'prep',
       '/static/images/places/prep-' || ((r.id + 1) % 2 + 1) || '.svg',
       'Illustration of the prep area at ' || r.name
FROM restaurants r
ON CONFLICT (restaurant_id, kind, url) DO NOTHING;

-- Good for mixed groups: restaurants with GF options AND a full menu for
-- people who aren't gluten-free. (Re-running just sets the same flags.)
UPDATE restaurants SET good_for_groups = true
WHERE id IN (1, 4, 5, 6, 9, 11, 12, 13);

-- Phone numbers for the call-ahead script. 555-0100 to 555-0199 are reserved
-- for fiction, so these can never ring a real business.
-- An UPDATE (not INSERT) because the restaurants already exist; "phone = ''"
-- means re-running never overwrites a number someone changed.
UPDATE restaurants SET phone = v.phone
FROM (VALUES
  (1,  '(801) 555-0101'), (2,  '(801) 555-0102'), (3,  '(801) 555-0103'), (4,  '(801) 555-0104'),
  (5,  '(801) 555-0105'), (6,  '(801) 555-0106'), (7,  '(801) 555-0107'), (8,  '(801) 555-0108'),
  (9,  '(801) 555-0109'), (11, '(303) 555-0111'), (12, '(303) 555-0112')
) AS v(id, phone)
WHERE restaurants.id = v.id AND restaurants.phone = '';

-- Safety profiles. Mile Marker Tacos (10) deliberately has none, to show the
-- "no safety information yet" state. Saltbox Thai (7) is out of date on purpose.
INSERT INTO safety_profiles (restaurant_id, verified_on, answers_updated_on, conflict_note) VALUES
  (1,  '2026-09-12', '2026-08-30', '2 reports disagree about whether GF pizzas bake in the same oven as regular pizzas.'),
  (2,  '2026-09-20', '2026-09-01', ''),
  (3,  '2026-07-02', '2026-06-15', 'Reports disagree about whether rice noodles are cooked in baskets also used for wheat noodles.'),
  (4,  '2026-09-05', '2026-08-20', '2 reports disagree about whether tortilla chips share a fryer with flour items.'),
  (5,  '2026-09-18', '2026-09-10', ''),
  (6,  '2026-09-08', '2026-08-25', ''),
  (7,  '2026-04-10', '2026-04-10', ''),
  (8,  '2026-09-25', '2026-09-15', ''),
  (9,  '2026-09-22', '2026-09-01', '1 report says GF buns were toasted on the shared grill; the restaurant says they use a separate pan.'),
  (11, '2026-08-14', '2026-08-01', ''),
  (12, '2026-09-02', '2026-08-18', 'Reports disagree about whether naan and tandoori dishes share the tandoor.'),
  -- Cluck Shack: reviewers have visited, but the restaurant hasn't answered (NULL).
  (13, '2026-09-15', NULL, '')
ON CONFLICT (restaurant_id) DO NOTHING;

-- Precaution badges. Labels name a practice, never a verdict; a "?" marks
-- the open question when reports conflict.
INSERT INTO precautions (restaurant_id, label, status, position) VALUES
  (1,  'Separate prep area',            'confirmed',  1),
  (1,  'Separate pizza trays',          'confirmed',  2),
  (1,  'Gloves changed',                'confirmed',  3),
  (1,  'Shared oven?',                  'conflict',   4),
  (2,  'Separate prep area',            'confirmed',  1),
  (2,  'Clean utensils',                'confirmed',  2),
  (2,  'Dedicated toaster',             'unverified', 3),
  (3,  'Rice noodles only',             'confirmed',  1),
  (3,  'Shared noodle baskets?',        'conflict',   2),
  (3,  'GF hoisin on request',          'unverified', 3),
  (4,  'Corn tortillas made in house',  'confirmed',  1),
  (4,  'Shared fryer?',                 'conflict',   2),
  (5,  'Staff trained on GF orders',    'confirmed',  1),
  (5,  'Gloves changed',                'confirmed',  2),
  (5,  'Separate prep area',            'confirmed',  3),
  (6,  'GF pasta in fresh water',       'confirmed',  1),
  (6,  'Separate pots',                 'confirmed',  2),
  (6,  'Separate colander',             'unverified', 3),
  (7,  'GF tamari on request',          'unverified', 1),
  (7,  'Clean wok between orders',      'unverified', 2),
  (8,  'Dedicated GF baking batch',     'confirmed',  1),
  (8,  'Dedicated toaster',             'confirmed',  2),
  (8,  'Separate prep area',            'confirmed',  3),
  (9,  'Dedicated fryer',               'confirmed',  1),
  (9,  'GF buns stored separately',     'confirmed',  2),
  (9,  'Shared grill?',                 'conflict',   3),
  (11, 'Sauces made without wheat',     'confirmed',  1),
  (11, 'Separate cutting boards',       'unverified', 2),
  (12, 'No wheat thickeners',           'confirmed',  1),
  (12, 'Shared tandoor?',               'conflict',   2),
  (13, 'Dedicated fryer',               'unverified', 1),
  (13, 'Separate prep area',            'unverified', 2)
ON CONFLICT (restaurant_id, label) DO NOTHING;

-- A few sample reviews so stars and reviewer tags have something to show.
-- Ids start at 1001 so they never collide with reviews posted in the app.
INSERT INTO reviews (id, meal_id, author_name, body, reviewer_type, rating, created_at) VALUES
  (1001, 1,  'Maya R.',   'Told the server I''m celiac. The cook came out, changed gloves and used a separate tray. Crust was crisp and chewy.', 'celiac',     5, '2026-09-12 19:10:00+00'),
  (1002, 1,  'Devon',     'Good pie. I asked about the oven and got two different answers from staff, so I''d call ahead.',                     'strict',     4, '2026-09-03 18:00:00+00'),
  (1003, 2,  'Priya K.',  'Pepperoni was great. GF crust is a little dense but tasty.',                                                          'sensitive',  4, '2026-08-28 20:30:00+00'),
  (1004, 3,  'Sam',       'Simple and fresh. They brought the dressing on the side when I asked.',                                               'celiac',     4, '2026-09-08 12:15:00+00'),
  (1005, 8,  'Lin',       'Broth was rich. One server said the noodle baskets are shared; another said they have separate ones.',                'celiac',     3, '2026-07-02 13:00:00+00'),
  (1006, 29, 'Jordan',    'Bun held together well. Staff knew exactly what I meant by cross-contact.',                                          'strict',     5, '2026-09-22 18:45:00+00'),
  (1007, 31, 'Alex T.',   'Fries come from a fryer used only for potatoes. They confirmed it without me pushing.',                              'celiac',     5, '2026-09-20 19:30:00+00'),
  (1008, 31, 'Beehive Burger Co.', 'Our fryer is used only for potatoes, and we filter the oil daily.',                                         'restaurant', 0, '2026-09-01 10:00:00+00'),
  (1009, 25, 'Riley',     'Fluffy pancakes, and the toaster for GF bread is labeled.',                                                          'sensitive',  5, '2026-09-25 09:20:00+00'),
  (1010, 39, 'Nadia',     'Delicious. I skipped the naan since they weren''t sure about the tandoor.',                                          'celiac',     4, '2026-09-02 19:00:00+00'),
  (1011, 46, 'Chris B.',  'Grilled chicken was juicy. They left the croutons off without me asking twice.',                                     'strict',     4, '2026-09-15 12:30:00+00'),
  (1012, 45, 'Morgan',    'Great fries, but the cashier confirmed they share oil with the breaded chicken, so I passed.',                       'celiac',     3, '2026-09-10 18:10:00+00')
ON CONFLICT (id) DO NOTHING;

-- Because we inserted explicit ids, move each id sequence past the highest id
-- so future inserts (e.g. from the app) don't collide.
SELECT setval(pg_get_serial_sequence('cities', 'id'),      (SELECT MAX(id) FROM cities));
SELECT setval(pg_get_serial_sequence('restaurants', 'id'), (SELECT MAX(id) FROM restaurants));
SELECT setval(pg_get_serial_sequence('meals', 'id'),       (SELECT MAX(id) FROM meals));
SELECT setval(pg_get_serial_sequence('reviews', 'id'),     (SELECT MAX(id) FROM reviews));
