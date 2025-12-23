-- Вставка данных в таблицу users
INSERT INTO users (id, email, password, name, status) VALUES
  (1, 'alexfit@gmail.com', '$2a$10$VeWhP15waRt1Vw9fNaVZ0OF/et0StBS7EGgrqrc9Ai9fmvAmUCo5C', 'Qwerty', 'banned'),
  (2, 'middif@gmail.com', '$2a$10$D.mo/.CF1/T3mFRSs9Bweug7IKYeiBnuGwB4hvGAaOmugKb8TLP2S', 'Виталя', 'user'),
  (3, 'bigbusness@gmail.com', '$2a$10$pfbOu9JkW0XVSKJZ8sNjjulj.L9DBNlDPjt7VYj3C.DDzpHvIes4e', 'Олег', 'user'),
  (4, 'basa@gmail.com', '$2a$10$MTyEapYm9Gmhcpkwa0D.0.K2F99RSCrLpvedxBTF4g9HCzHO7kc9W', 'Иван', 'user'),
  (5, 'elena1915@mail.ru', '$2a$10$.yKfarKgUtd5MoUNAWkCy.0xBPNIdiud8/Hll/6viVHi0.87QCDmm', 'Елена', 'user'),
  (6, 'top228@gmail.com', '$2a$10$zyvbRgnCtm6SRE3WvY77CuGMS4DlYVW7EM/lN7kRrf0qvmE7GYC/q', 'Виктор', 'user'),
  (7, 'greenfog@gmail.com', '$2a$10$JiptWmkN2jHDdXxbPFM.r.9g2Cq4X1WK86STS0n5xzmp1.o8t9XqG', 'Юля', 'user'),
  (8, '5rublei@gmail.com', '$2a$10$zcI.ICYzGFkFIffFwTZDRernoJrT30ljmWthnKwGtB0iN5I4M8BPK', 'Снежана', 'user'),
  (9, 'a@a.com', '$2a$10$wFALUJX5E0asKBAY5sy30OtMKWEY/HsCcw9aE0nu62mg/7f0GNLYW', 'a', 'user'),
  (10, 'gw@gmail.com', '$2a$10$WFhmXqLBo004Qh9NsAYXTOUvPCPIEHkJQkPO4uyvb9eUX.EuOUbYq', 'administator', 'user'),
  (11, 'qwerty@gmail.com', '$2a$10$HnFmoGuJtbkPfvbsb24oMODQ8YJZgKDOq15QqWIQqS1.6/qoRk3P.', 'qwerty', 'admin'),
  (12, 'brintol@gmail.com', '$2a$10$z86VxT3scZfxk3sezACUte6HNpsXJhK8nucxAzGImHPubC.2v1L2K', 'baobab', 'user'),
  (13, 'yguhkbhi@gmail.com', '$2a$10$pmhfKqV3YYZxbTZ7H4iIIOBXTrRIglGj.gP0uLZwTaK3FUKKSpgKy', 'fgjhjig', 'user'),
  (15, 'dasha@gmail.com', '$2a$10$UF0qjtLem7gezTwXae0db.gWy.r0bG4dYZphdOEMghrYnNP02M9ie', 'Dasha', 'user'),
  (18, 'basa1@gmail.com', '$2a$10$dUJi1hNpQH.Gacuh65bLyuYdIGxQ.a9Gbf4Hy45kf6rVbeIu9xXGa', 'ffff', 'user'),
  (20, 'qwerty123@gmail.com', '$2a$10$hSm4QHqT9DWFZmfOIqm12OWc3ktvEvq8PT2SJ4Mt6DCCnYShN.K.W', 'qwerty123', 'user'),
  (32, 'perelajkovlad@gmail.com', '$2a$10$U707TnmPBhO1EZ2ZSrgZT.c6W0KzAeiJlbcUovJsPCyWRh/cyrJdK', 'Vlad', 'admin');

-- Сброс последовательности для users
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 0) + 1);

-- Вставка данных в таблицу categories
INSERT INTO categories (category_id, name, description) VALUES
  (1, 'Протеин', 'Белковые добавки для роста и восстановления мышц'),
  (2, 'Креатин', 'Добавки для увеличения силы и выносливости'),
  (3, 'BCAA', 'Аминокислоты с разветвленной цепью для восстановления'),
  (4, 'Батончики', 'Протеиновые и низкоуглеводные батончики'),
  (5, 'Витамины', 'Витаминно-минеральные комплексы'),
  (6, 'Предтреники', 'Предтренировочные комплексы для энергии');

-- Сброс последовательности для categories
SELECT setval('categories_category_id_seq', COALESCE((SELECT MAX(category_id) FROM categories), 0) + 1);

-- Вставка данных в таблицу producers
INSERT INTO producers (producer_id, name, country) VALUES
  (1, 'Optimum Nutrition', 'USA'),
  (2, 'Maxler', 'Germany'),
  (3, 'Bombbar', 'Russia'),
  (4, 'Power System', 'Germany'),
  (5, 'MyProtein', 'UK'),
  (6, 'GeneticLab', 'Russia'),
  (7, '1WIN', 'Russia'),
  (8, 'Prime Kraft', 'Russia');

-- Сброс последовательности для producers
SELECT setval('producers_producer_id_seq', COALESCE((SELECT MAX(producer_id) FROM producers), 0) + 1);

-- Вставка данных в таблицу products
INSERT INTO products (product_id, name, description, price, stock_qty, image_url) VALUES
  (1, 'Optimum Nutrition Gold Standard 100% Whey 600g', 'Сывороточный протеин 600 г', 99.90, 14, 'Optimum_Nutrition_Gold_Standard_100%_Whey_600g.webp'),
  (2, 'Optimum Nutrition Gold Standard 100% Whey 900g', 'Сывороточный протеин 900 г ', 159.90, 14, 'Optimum_Nutrition_Gold_Standard_100%_Whey_900g.webp'),
  (3, 'Optimum Nutrition Gold Standard 100% Whey 4500g', 'Сывороточный протеин 4500 г', 499.90, 1, 'Optimum_Nutrition_Gold_Standard_100%_Whey_4500g.webp'),
  (4, 'Maxler 100% Golden Whey 900g', 'Сывороточный протеин 900 г', 129.90, 17, 'Maxler_100%_Golden_Whey_900g.webp'),
  (5, 'Maxler 100% Golden Whey 2270g', 'Сывороточный протеин 2270 г', 259.90, 19, 'Maxler_100%_Golden_Whey_2270g.webp'),
  (6, 'GeneticLab Whey Pro 150g', 'Изолят сывороточного протеина 150 г', 29.90, 12, 'GeneticLab_Whey_Pro_150g.webp'),
  (7, 'GeneticLab Whey Pro 2100g', 'Изолят сывороточного протеина 2100 г', 239.90, 13, 'GeneticLab_Whey_Pro_2100g.webp'),
  (8, '1WIN WHEY PROTEIN 450g', 'Сывороточный протеин 450 г', 69.90, 18, '1WIN_WHEY_PROTEIN_450g.webp'),
  (9, '1WIN WHEY PROTEIN 900g', 'Сывороточный протеин 900 г', 119.90, 9, '1WIN_WHEY_PROTEIN_900g.webp'),
  (10, 'Optimum Nutrition Micronized Creatine 250g', 'Креатин моногидрат 250 г', 49.90, 26, 'Optimum_Nutrition_Micronized_Creatine_250g.webp'),
  (11, 'Optimum Nutrition Micronized Creatine 2500g', 'Креатин моногидрат 2500 г', 449.90, 24, 'Optimum_Nutrition_Micronized_Creatine_2500g.webp'),
  (12, 'Maxler Creatine Monohydrate 300g', 'Креатин моногидрат 300 г', 59.90, 26, 'Maxler_Creatine_Monohydrate_300g.webp'),
  (13, 'Maxler Creatine Monohydrate 500g', 'Креатин моногидрат 500 г', 79.90, 20, 'Maxler_Creatine_Monohydrate_500g.webp'),
  (14, 'Power System Creatine 650g', 'Креатин моногидрат 650 г', 99.90, 25, 'Power_System_Creatine_650g.webp'),
  (15, 'Power System Creatine 3000g', 'Креатин моногидрат 3000 г', 299.90, 6, 'Power_System_Creatine_3000g.webp'),
  (16, 'Prime Kraft Creatine Monohydrate 200g', 'Креатин моногидрат 200 г', 39.90, 22, 'Prime_Kraft_Creatine_Monohydrate_200g.webp'),
  (17, 'Prime Kraft Creatine Monohydrate 500g', 'Креатин моногидрат 500 г', 79.90, 18, 'Prime_Kraft_Creatine_Monohydrate_500g.webp'),
  (18, 'MyProtein BCAA 250g', 'BCAA 2:1:1 250 г', 69.90, 35, 'MyProtein_BCAA_250g.webp'),
  (19, 'MyProtein BCAA 500g', 'BCAA 2:1:1 500 г', 119.90, 24, 'MyProtein_BCAA_500g.webp'),
  (20, 'GeneticLab BCAA Pro 250g', 'BCAA 4:1:1 250 г', 79.90, 20, 'GeneticLab_BCAA_Pro_250g.webp'),
  (21, 'GeneticLab BCAA Pro 500g', 'BCAA 4:1:1 500 г', 129.90, 15, 'GeneticLab_BCAA_Pro_500g.webp'),
  (22, '1WIN BCAA 200g', 'BCAA с электролитами 200 г', 59.90, 25, '1WIN_BCAA_200g.webp'),
  (23, '1WIN BCAA 360g', 'BCAA с электролитами 360 г', 99.90, 20, '1WIN_BCAA_360g.webp'),
  (24, 'Bombbar Protein Bar 60g', 'Протеиновый батончик 60 г', 4.90, 100, 'Bombbar_Protein_Bar_60g.webp'),
  (25, 'Bombbar Wafer 45g', 'Вафельный батончик 45 г', 4.50, 55, 'Bombbar_Wafer_45g.webp'),
  (26, 'Maxler Protein Bar 65g', 'Протеиновый батончик 65 г', 5.90, 56, 'Maxler_Protein_Bar_65g.webp'),
  (27, '1WIN Pre-Workout 210g', 'Предтренировочный комплекс 210 г', 59.90, 13, '1WIN_Pre-Workout_210g.webp'),
  (28, '1WIN Pre-Workout 400g', 'Предтренировочный комплекс 400 г', 99.90, 10, '1WIN_Pre-Workout_400g.webp'),
  (29, 'GeneticLab Pre-Workout 200g', 'Предтренировочный комплекс 200 г', 55.90, 11, 'GeneticLab_Pre-Workout_200g.webp'),
  (30, 'Maxler Vita Men 90 capsules', 'Витамины для мужчин 90 капсул', 59.90, 39, 'Maxler_Vita_Men_90_capsules.webp'),
  (31, 'Maxler Vita Women 90 capsules', 'Витамины для женщин 90 капсул', 59.90, 30, 'Maxler_Vita_Women_90_capsules.webp'),
  (32, 'Optimum Nutrition Opti-Men 120 capsules', 'Мультивитамины для мужчин 120 капсул', 69.90, 35, 'Optimum_Nutrition_Opti-Men_120_capsules.webp'),
  (33, 'Optimum Nutrition Opti-Women 120 capsules', 'Мультивитаминыдля женщин 120 капсул', 69.90, 25, 'Optimum_Nutrition_Opti-Women_120_capsules.webp');

-- Сброс последовательности для products
SELECT setval('products_product_id_seq', COALESCE((SELECT MAX(product_id) FROM products), 0) + 1);

-- Вставка данных в таблицу products_categories
INSERT INTO products_categories (product_id, category_id) VALUES
  (1, 1),
  (2, 1),
  (3, 1),
  (4, 1),
  (5, 1),
  (6, 1),
  (7, 1),
  (8, 1),
  (9, 1),
  (10, 2),
  (11, 2),
  (12, 2),
  (13, 2),
  (14, 2),
  (15, 2),
  (16, 2),
  (17, 2),
  (18, 3),
  (19, 3),
  (20, 3),
  (21, 3),
  (22, 3),
  (23, 3),
  (24, 4),
  (25, 4),
  (26, 4),
  (27, 6),
  (28, 6),
  (29, 6),
  (30, 5),
  (31, 5),
  (32, 5),
  (33, 5);

-- Вставка данных в таблицу products_producers
INSERT INTO products_producers (product_id, producer_id) VALUES
  (1, 1),
  (2, 1),
  (3, 1),
  (4, 2),
  (5, 2),
  (6, 6),
  (7, 6),
  (8, 7),
  (9, 7),
  (10, 1),
  (11, 1),
  (12, 2),
  (13, 2),
  (14, 4),
  (15, 4),
  (16, 8),
  (17, 8),
  (18, 5),
  (19, 5),
  (20, 6),
  (21, 6),
  (22, 7),
  (23, 7),
  (24, 3),
  (25, 3),
  (26, 2),
  (27, 7),
  (28, 7),
  (29, 6),
  (30, 2),
  (31, 2),
  (32, 1),
  (33, 1);

-- Вставка данных в таблицу carts
INSERT INTO carts (cart_id, user_id) VALUES
  (1, 11),
  (2, 12),
  (3, 13),
  (5, 15),
  (8, 18),
  (10, 20),
  (13, 32);

-- Сброс последовательности для carts
SELECT setval('carts_cart_id_seq', COALESCE((SELECT MAX(cart_id) FROM carts), 0) + 1);

-- Вставка данных в таблицу cart_items
INSERT INTO cart_items (cart_item_id, cart_id, product_id, quantity) VALUES
  (12, 5, 12, 1),
  (65, 10, 2, 1),
  (66, 10, 1, 1),
  (67, 10, 29, 1),
  (68, 13, 2, 1);

-- Сброс последовательности для cart_items
SELECT setval('cart_items_cart_item_id_seq', COALESCE((SELECT MAX(cart_item_id) FROM cart_items), 0) + 1);

-- Вставка данных в таблицу orders
INSERT INTO orders (order_id, user_id, created_at, total_amount) VALUES
  (5, 11, '2025-11-07 16:09:36.662115+03', 17.70),
  (6, 11, '2025-11-07 16:10:44.913398+03', 549.70),
  (7, 11, '2025-11-07 17:20:55.483787+03', 319.60),
  (8, 11, '2025-11-08 15:55:15.266624+03', 299.70),
  (9, 11, '2025-11-09 00:02:07.497018+03', 99.90),
  (11, 18, '2025-11-10 15:39:52.002896+03', 149.70),
  (13, 20, '2025-11-12 14:41:56.6649+03', 159.90),
  (14, 20, '2025-11-12 15:37:13.921521+03', 159.90),
  (16, 20, '2025-11-27 18:42:53.62086+03', 2770.70),
  (17, 20, '2025-12-10 00:22:45.699523+03', 159.90),
  (21, 20, '2025-12-12 22:40:44.295521+03', 259.90);

-- Сброс последовательности для orders
SELECT setval('orders_order_id_seq', COALESCE((SELECT MAX(order_id) FROM orders), 0) + 1);

-- Вставка данных в таблицу order_items
INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
  (5, 26, 3, 5.90),
  (6, 7, 2, 239.90),
  (6, 8, 1, 69.90),
  (7, 19, 1, 119.90),
  (7, 30, 1, 59.90),
  (7, 33, 2, 69.90),
  (8, 1, 3, 99.90),
  (9, 1, 1, 99.90),
  (11, 10, 3, 49.90),
  (13, 2, 1, 159.90),
  (14, 2, 1, 159.90),
  (16, 26, 1, 5.90),
  (16, 12, 2, 59.90),
  (16, 2, 1, 159.90),
  (16, 4, 1, 129.90),
  (16, 29, 1, 55.90),
  (16, 1, 3, 99.90),
  (16, 3, 4, 499.90),
  (17, 2, 1, 159.90),
  (21, 5, 1, 259.90);