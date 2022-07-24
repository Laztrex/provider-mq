import logging


# from model.initial import Initializer

def start_predict(data: list) -> int:
    """
    Пример cpu-bounded функции
    :return:
    """
    # print(Initializer.DF_SHARED_VALUE.df)

    sum_predict = 0
    for num in data:
        logging.info(f'getting num = {num}')
        sum_predict += sum(list(fibb_gen(num)))

    res = len(str(2 ** sum_predict))

    logging.info('DONE')

    return res


def fibb_gen(n):
    a, b = 0, 1
    for _ in range(n):
        yield a
        a, b = b, a + b
