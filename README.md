Реализовать модель обработки данных в виде пайплайна, состоящего из следующих этапов

1. Подача на вход пакетов данных. Пакет данных = слайсу случайных целых чисел из 10 элементов. Новый пакет подается каждые N мс (N задается в виде env переменной)   

2. Обработка пакетов: нахождение 3-х наибольших чисел в пакете. Вход: слайс int из 10 элементов, выход: слайс из 3-х элементов. Обработка пакетов должна производиться M воркерами (M задается в виде env переменной)

3. Аккумулятор: суммирование чисел обработанных пакетов, полученных на предыдущем этапе, и запись в единую переменную int

4. Публикатор: вывод на консоль текущего значения аккумулятора каждые K секунд (K задается в виде env переменной)

Пример:
вход: {1, 9, 6, 4, 4, 5, 7, 8, 0, 1}
обработка: {9, 7, 8}
аккумулятор: 9+7+8=24
публикатор: 24


