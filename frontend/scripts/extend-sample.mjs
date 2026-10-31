// Grows src/data/sample.json so testers see a realistic amount of content.
// Everything here is fictional. Run with: node scripts/extend-sample.mjs
//
// The first 13 restaurants and 59 meals come from the Go app's seed. This script
// keeps those and rebuilds everything after them, so it's safe to re-run, and a
// seeded random generator means re-running gives the same result.
import { readFileSync, writeFileSync } from "node:fs";

const PATH = new URL("../src/data/sample.json", import.meta.url);
const data = JSON.parse(readFileSync(PATH, "utf8"));
const PROVO = 1, SLC = 2, DENVER = 3;
const AREA = { [PROVO]: "801", [SLC]: "801", [DENVER]: "303" };
const TOWN = { [PROVO]: "Provo, UT", [SLC]: "Salt Lake City, UT", [DENVER]: "Denver, CO" };

let seed = 581;
const rand = () => {
  seed = (seed + 0x6d2b79f5) | 0;
  let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
  t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
};
const pick = (list) => list[Math.floor(rand() * list.length)];

// [name, description, image, status, note, food type]
const gf = (name, description, photo, foodType = "entree", gfNote = "") => ({ name, description, photo, gfStatus: "gf", gfNote, foodType });
const cross = (name, description, photo, gfNote, foodType = "entree") => ({ name, description, photo, gfStatus: "cross_contact", gfNote, foodType });
const contains = (name, description, photo, gfNote, foodType = "entree") => ({ name, description, photo, gfStatus: "contains", gfNote, foodType });

const MENUS = {
  Italian: [
    gf("GF Penne Pomodoro", "Corn-rice penne with slow-cooked tomato and basil.", "pasta", "grain", "Boiled in its own pot."),
    gf("Risotto alla Milanese", "Creamy saffron risotto with parmesan.", "bowl", "grain"),
    gf("Caprese Salad", "Fresh mozzarella, tomato and basil.", "salad", "salad"),
    gf("Lemon Herb Branzino", "Roasted whole fish with lemon and capers.", "fish"),
    cross("GF Margherita Pizza", "Cauliflower crust, tomato and mozzarella.", "pizza", "Baked in the shared oven.", "grain"),
    contains("Fettuccine Alfredo", "Wheat fettuccine in a parmesan cream sauce.", "pasta", "Made with wheat pasta.", "grain"),
    contains("Garlic Breadsticks", "Warm, buttery and made in house.", "bread", "Made with wheat flour.", "grain"),
  ],
  Mexican: [
    gf("Carnitas Corn Tacos", "Slow-braised pork on house corn tortillas.", "taco"),
    gf("Chicken Tinga Bowl", "Smoky chicken over rice and black beans.", "bowl", "grain"),
    gf("Pozole Rojo", "Hominy and pork in a red chile broth.", "soup"),
    gf("Ensalada de Nopal", "Cactus, tomato, cotija and lime.", "salad", "salad"),
    cross("House Tortilla Chips", "Crisp corn chips with salsa roja.", "fries", "Fried in a fryer shared with flour items."),
    contains("Flour Tortilla Burrito", "Carne asada, rice and beans in a flour tortilla.", "taco", "Made with a wheat tortilla."),
    contains("Churros", "Cinnamon sugar with chocolate sauce.", "dessert", "Made with wheat flour.", "dessert"),
  ],
  Thai: [
    gf("Pad Thai with Shrimp", "Rice noodles, tamarind, peanuts and lime.", "noodles", "grain", "Made with GF tamari."),
    gf("Green Curry with Chicken", "Coconut curry with Thai eggplant and basil.", "curry"),
    gf("Larb Gai Salad", "Minced chicken, herbs and toasted rice.", "salad", "salad"),
    gf("Tom Kha Soup", "Coconut soup with galangal and mushrooms.", "soup"),
    cross("Crispy Basil Tofu", "Rice-flour tofu with holy basil.", "bowl", "Fried in shared oil."),
    contains("Crispy Spring Rolls", "Vegetable rolls with sweet chili sauce.", "rolls", "Wrappers contain wheat."),
    contains("Beef with Oyster Sauce", "Stir-fried beef and broccoli.", "steak", "The sauce contains wheat."),
  ],
  American: [
    gf("Blue Plate Meatloaf", "Oat-bound meatloaf with mashed potatoes.", "steak"),
    gf("Grilled Chicken Plate", "Herb chicken with roasted vegetables.", "chicken"),
    gf("Cobb Salad", "Chicken, bacon, egg, avocado and blue cheese.", "salad", "salad"),
    gf("GF Pancake Stack", "Buttermilk-style pancakes with maple syrup.", "breakfast", "grain", "Cooked on a separate griddle."),
    cross("Hand-Cut Fries", "Crispy and salted.", "fries", "Fried in a fryer shared with breaded items."),
    contains("Buttermilk Biscuits and Gravy", "Flaky biscuits with sausage gravy.", "bread", "Made with wheat flour.", "grain"),
    contains("Patty Melt", "Beef and onions on griddled rye.", "sandwich", "Served on rye bread."),
  ],
  Chinese: [
    gf("Steamed Ginger Fish", "Whole fish with ginger, scallion and GF tamari.", "fish"),
    gf("Egg Fried Rice", "Wok-fried rice with peas and GF tamari.", "bowl", "grain"),
    gf("Rice Noodle Stir-Fry", "Wide rice noodles with chicken and greens.", "noodles", "grain", "Cooked in a separate wok."),
    gf("Hot and Sour Soup", "Thickened with cornstarch.", "soup"),
    cross("Salt and Pepper Shrimp", "Lightly dusted in rice flour.", "fish", "Fried in shared oil."),
    contains("Pork Potstickers", "Pan-fried dumplings.", "rolls", "Wrappers contain wheat."),
    contains("Kung Pao Chicken", "Peanuts, chiles and soy sauce.", "chicken", "Made with regular soy sauce."),
  ],
  Bakery: [
    gf("GF Snickerdoodle", "Soft, cinnamon-rolled and baked daily.", "dessert", "dessert", "Baked in a dedicated GF batch."),
    gf("Almond Flour Brownie", "Dense and fudgy with sea salt.", "dessert", "dessert"),
    gf("Lemon Polenta Cake", "Bright, moist and lightly sweet.", "dessert", "dessert"),
    gf("GF Sandwich Loaf", "Sliced to order.", "bread", "grain"),
    gf("Quinoa Breakfast Bowl", "Warm quinoa with berries and honey.", "bowl", "grain"),
    cross("Coconut Macaroons", "Crisp outside and chewy inside.", "dessert", "Baked in an oven shared with wheat goods.", "dessert"),
    contains("Butter Croissant", "Flaky and laminated in house.", "bread", "Made with wheat flour.", "grain"),
  ],
};
const DESSERT = gf("Flourless Torte", "Dark chocolate torte with raspberry.", "dessert", "dessert");

// id, city, name, street, cuisine, description, good for groups, has a GF dessert
const NEW_RESTAURANTS = [
  [14, PROVO, "Nonna Bea's Table", "61 Olive Row", "Italian", "Family-style Italian with a marked GF pasta list.", true, true],
  [15, PROVO, "Casa Verde Cantina", "204 Mesquite St", "Mexican", "Corn tortillas pressed in house and a big group patio.", true, false],
  [16, PROVO, "Golden Basil Thai", "718 Cedar Fork Rd", "Thai", "A small family kitchen cooking with GF tamari on request.", false, false],
  [17, PROVO, "Lemongrass Alley", "33 Paper Mill Ln", "Thai", "Street-style Thai plates and long shared tables.", true, false],
  [18, PROVO, "Honeycrumb Bakery", "9 Sugarbeet Way", "Bakery", "Daily GF batches of cookies, cakes and bread.", false, true],
  [19, SLC, "Rio Azul Cocina", "480 Copper Ave", "Mexican", "Coastal Mexican with plenty for the whole table.", true, false],
  [20, SLC, "Agave & Ember", "1275 Foundry Pl", "Mexican", "Wood-fired Mexican plates and corn-only tacos.", false, false],
  [21, SLC, "Trattoria Alba", "66 Marble Row", "Italian", "Handmade pasta, with GF pasta cooked separately.", true, true],
  [22, SLC, "Orchid Room Thai", "340 Lantern Sq", "Thai", "Regional Thai curries in a roomy dining room.", true, false],
  [23, SLC, "Jade Lantern", "802 Willow Bend", "Chinese", "Cantonese classics, with a GF wok on request.", true, false],
  [24, DENVER, "High Plains Diner", "2100 Prairie Rd", "American", "All-day diner plates and a separate GF griddle.", true, true],
  [25, DENVER, "Piazza Verde", "418 Terrace Ct", "Italian", "Neighborhood Italian with a dedicated GF dessert case.", true, true],
  [26, DENVER, "Mesa Pasta Co.", "95 Red Clay St", "Italian", "A counter-service pasta shop with GF options.", false, false],
  [27, DENVER, "Cielo Taqueria", "1610 Altitude Ave", "Mexican", "Late-night tacos on corn tortillas.", true, false],
  [28, DENVER, "Basil & Lime", "57 Pinecone Way", "Thai", "Fresh Thai bowls and curries made to order.", false, false],
  [29, DENVER, "Bangkok Junction", "3300 Railyard Blvd", "Thai", "Big-table Thai near the old rail yard.", true, false],
  [30, DENVER, "Snowline Sweets", "12 Glacier Ln", "Bakery", "A small dessert shop with several GF bakes.", false, true],
];

const restaurants = data.restaurants.filter((r) => r.id <= 13);
const meals = data.meals.filter((m) => m.id <= 59);
let mealId = 60;
for (const [id, cityId, name, street, cuisine, description, goodForGroups, hasDessert] of NEW_RESTAURANTS) {
  restaurants.push({
    id, cityId, name, cuisine, description, goodForGroups,
    address: `${street}, ${TOWN[cityId]}`,
    // 555-0100 to 555-0199 are reserved for fiction, so these can't ring anyone.
    phone: id === 28 ? "" : `(${AREA[cityId]}) 555-01${id + 30}`,
  });
  const menu = MENUS[cuisine];
  const start = id % 3;
  // Five or six dishes, always keeping one gluten-containing item for the group.
  const chosen = [...menu.slice(start, start + (id % 2 ? 4 : 5)), menu[menu.length - 1]];
  if (hasDessert && cuisine !== "Bakery") chosen.splice(1, 0, DESSERT);
  for (const dish of chosen) {
    meals.push({ id: mealId++, restaurantId: id, ...dish, gfNote: dish.gfNote, photo: `/images/meals/${dish.photo}.svg` });
  }
}

// Safety profiles, keeping every state testers should run into.
const CONFLICTS = {
  15: "2 reports disagree about whether the tortilla chips share a fryer with flour items.",
  22: "Reports disagree about whether the wok is cleaned between orders.",
  27: "1 report says tacos were warmed on a shared grill; the restaurant says they use a separate pan.",
};
const PRACTICES = ["Separate prep area", "Clean utensils", "Gloves changed", "Dedicated fryer", "Separate cookware"];
const safetyProfiles = data.safetyProfiles.filter((s) => s.restaurantId <= 13);
const precautions = data.precautions.filter((p) => p.restaurantId <= 13);
for (const [id] of NEW_RESTAURANTS) {
  if (id === 30) continue; // no profile: the "no safety information" state
  const stale = id === 26; // out of date on purpose
  const unanswered = id === 18;
  const day = String((id % 27) + 1).padStart(2, "0");
  safetyProfiles.push({
    restaurantId: id,
    verifiedOn: stale ? "2026-03-02" : `2026-${id % 2 ? "09" : "10"}-${day}`,
    answersUpdatedOn: unanswered ? null : stale ? "2026-03-02" : `2026-09-${day}`,
    conflictNote: CONFLICTS[id] ?? "",
  });
  const start = id % PRACTICES.length;
  for (let i = 0; i < 3; i++) {
    const label = PRACTICES[(start + i) % PRACTICES.length];
    let status = "confirmed";
    if (stale || (i === 2 && id % 3 === 0)) status = "unverified";
    if (CONFLICTS[id] && i === 2) status = "conflict";
    precautions.push({ restaurantId: id, label, status });
  }
}

// Restaurant photos: the restaurant itself, its prep area and its menu. Food
// photos live on meals instead, because those come from diners.
const photos = restaurants.flatMap((r) => [
  { restaurantId: r.id, kind: "restaurant", url: `/images/places/storefront-${(r.id % 2) + 1}.svg`, alt: `Illustration of the outside of ${r.name}` },
  { restaurantId: r.id, kind: "prep", url: `/images/places/prep-${((r.id + 1) % 2) + 1}.svg`, alt: `Illustration of the prep area at ${r.name}` },
  { restaurantId: r.id, kind: "menu", url: `/images/places/menu-${(r.id % 2) + 1}.svg`, alt: `Illustration of a menu from ${r.name}`, addedOn: `2026-${r.id % 2 ? "09" : "10"}-${String((r.id % 20) + 1).padStart(2, "0")}` },
]);

// Reviews, with a mix of tones so testers aren't reading all five-star praise.
// Copy rule: reviews say what happened, and don't call a meal "safe".
const AUTHORS = ["Maya R.", "Devon", "Priya K.", "Sam", "Lin", "Jordan T.", "Aisha", "Marco", "Beth W.", "Noah", "Elena", "Theo B.", "Grace", "Hannah L.", "Omar", "Kira", "Ben S.", "Tasha", "Wes", "Nina P.", "Caleb", "Rosa", "Ivy M.", "Luke", "Dana", "Yuki", "Peter H.", "Amara", "Quinn", "Sofia"];
const TYPES = ["celiac", "celiac", "celiac", "strict", "strict", "sensitive"];
const LINES = {
  5: [
    "Best {meal} I've had since going gluten-free. The server checked with the kitchen without me asking.",
    "The {meal} was excellent, and they changed gloves right where I could see it. I'd order it again.",
    "I brought my whole family and everyone liked it. The {meal} didn't feel like a compromise.",
    "Really good {meal}. Staff knew exactly which sauces had gluten.",
    "Felt great afterward. I've ordered the {meal} twice now and it's been consistent.",
  ],
  4: [
    "Tasty {meal}. It took a while to come out, but they were careful about my order.",
    "Good flavor on the {meal}, though the portion was on the small side. Felt fine afterward.",
    "Solid {meal}. I had to explain what I needed, but the manager came out to confirm.",
    "I liked the {meal}. Ask for it without the garnish if you want to be extra careful.",
    "Nice {meal}. The menu marks GF clearly, which made ordering easy.",
  ],
  3: [
    "The {meal} was fine but nothing special. Staff seemed unsure at first and had to go ask.",
    "Decent {meal}. It was a little bland, but I felt okay afterward.",
    "The {meal} was okay. I'd call ahead next time, since the person I talked to wasn't sure about the prep.",
    "Average {meal}. Worth it if you're in the area, but I wouldn't make a special trip.",
  ],
  2: [
    "Tasted good, but I felt off afterward. It could be unrelated, but I'd ask about shared equipment.",
    "The {meal} was dry, and the server didn't seem to know what gluten-free meant. I'd ask a lot of questions here.",
    "My {meal} came out with a garnish I'd asked them to leave off. They fixed it, but it shook my confidence.",
  ],
};
const RESTAURANT_LINES = [
  "We prepare this at a separate station. Call ahead and we can walk you through how we make it.",
  "We updated this recipe to use GF tamari. Ask for the gluten-free menu when you sit down.",
  "We use clean pans for gluten-free orders. Please tell your server about your needs when you order.",
];
const TARGET_NEW_REVIEWS = 106;
// These stay review-free on purpose, for the "No reviews yet" state.
const NO_REVIEWS = new Set([2, 4, 20, 28]);
const eligible = meals.filter((m) => m.gfStatus !== "contains" && !NO_REVIEWS.has(m.restaurantId));
const reviews = data.reviews.filter((r) => r.id < 2000);
const start = Date.UTC(2026, 7, 1);
const span = Date.UTC(2026, 9, 28) - start;
const nameOf = (id) => restaurants.find((r) => r.id === id).name;
for (let i = 0; i < TARGET_NEW_REVIEWS; i++) {
  const meal = pick(eligible);
  const createdAt = new Date(start + rand() * span).toISOString();
  if (i % 25 === 12) {
    reviews.push({ id: 2000 + i, mealId: meal.id, author: nameOf(meal.restaurantId), body: pick(RESTAURANT_LINES), reviewerType: "restaurant", rating: 0, createdAt });
    continue;
  }
  const roll = rand();
  const rating = roll < 0.35 ? 5 : roll < 0.7 ? 4 : roll < 0.9 ? 3 : 2;
  const body = pick(LINES[rating]).replace("{meal}", meal.name.replace(/^GF /, "").toLowerCase());
  reviews.push({ id: 2000 + i, mealId: meal.id, author: pick(AUTHORS), body, reviewerType: pick(TYPES), rating, createdAt });
}

writeFileSync(PATH, `${JSON.stringify({ cities: data.cities, restaurants, meals, reviews, safetyProfiles, precautions, photos }, null, 1)}\n`);
console.log({ restaurants: restaurants.length, meals: meals.length, reviews: reviews.length, profiles: safetyProfiles.length, photos: photos.length });
