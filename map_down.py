import argparse
import asyncio
import json
import os

import aiohttp
import js2py
import requests
from tqdm import tqdm

# 重试次数常量
RETRY_TIMES = 3

# 设置图片下载 URL 模板
base_url = "https://palworld.gg/images/tiles/{z}/{x}/{y}.png"

# 本地保存文件的根目录
save_dir = "./map"

# 为图片请求设置 headers
headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
}

# 根据z的值设置x和y的最大范围
z_to_range = {
    0: (0, 0),  # z=0, x/y max=0
    1: (1, 1),  # z=1, x/y max=1
    2: (3, 3),  # z=2, x/y max=3
    3: (7, 7),  # z=3, x/y max=7
    4: (15, 15),  # z=4, x/y max=15
    5: (31, 31),  # z=5, x/y max=31
    6: (63, 63),  # z=6, x/y max=63
}


async def download_image(session, url, file_path, custom_headers, progress_bar, redown):
    # 检查本地文件是否已经存在
    if os.path.exists(file_path) and not redown:
        progress_bar.update(1)
        return

    # 尝试下载图片
    attempt = 0
    success = False
    while attempt < RETRY_TIMES:
        try:
            async with session.get(url, headers=custom_headers) as response:
                # 检查响应状态码
                if response.status == 200:
                    # 将图片内容保存到本地
                    with open(file_path, "wb") as f:
                        f.write(await response.read())
                    success = True
                    break  # 成功下载后退出重试循环
                elif response.status == 404:
                    print(f"Skipped {url} - Not Found (404)")
                    break  # 404 不重试
                elif response.status == 403:
                    print(f"Skipped {url} - Forbidden (403)")
                    break  # 403 不重试
                else:
                    print(f"Failed to download {url} (status code: {response.status})")
                    break  # 其他错误码不重试
        except Exception as e:
            print(f"Error downloading {url}: {e}")
        attempt += 1
        if attempt < RETRY_TIMES:
            print(f"Retrying ({attempt}/{RETRY_TIMES}) for {url}")
            await asyncio.sleep(1)  # 重试前等待 1 秒

    # 更新进度条
    if success:
        progress_bar.update(1)
    else:
        print(f"Failed to download {url} after {RETRY_TIMES} attempts")
        progress_bar.update(1)  # 即使失败也更新进度条


async def download_images_async(redown=False):
    # 计算总图片数
    total_images = sum(
        (x_max + 1) * (y_max + 1) for x_max, y_max in z_to_range.values()
    )

    # 初始化 tqdm 进度条
    progress_bar = tqdm(total=total_images, desc="Downloading images", unit="img")

    async with aiohttp.ClientSession() as session:
        tasks = []
        # 遍历 z, x, y 的范围
        for z, (x_max, y_max) in z_to_range.items():
            for x in range(0, x_max + 1):
                for y in range(0, y_max + 1):
                    # 构造图片的 URL
                    url = base_url.format(z=z, x=x, y=y)
                    # 构造本地保存的文件路径
                    save_path = os.path.join(save_dir, str(z), str(x))
                    file_name = f"{y}.png"
                    file_path = os.path.join(save_path, file_name)

                    # 如果文件路径不存在则创建
                    os.makedirs(save_path, exist_ok=True)

                    # 创建任务
                    task = download_image(
                        session, url, file_path, headers, progress_bar, redown
                    )
                    tasks.append(task)

        # 等待所有任务完成
        await asyncio.gather(*tasks)

    progress_bar.close()


def parse_js_file():
    # 下载JavaScript文件
    url = "https://paldb.cc/js/map_data_cn.js"
    response = requests.get(url, timeout=10)
    with open("map_data_cn.js", "wb") as file:
        file.write(response.content)

    # 读取JavaScript文件内容
    with open("map_data_cn.js", "r", encoding="utf-8") as file:
        js_content = file.read()

    # 使用js2py执行JavaScript代码
    context = js2py.EvalJs()
    context.execute(js_content)

    # 提取变量并转换为Python字典
    fixed_dungeon_obj = context.fixedDungeon.to_dict()
    fixed_dungeon = list(fixed_dungeon_obj.values())

    # 初始化结果字典
    result = {"boss_tower": [], "fast_travel": []}

    # 过滤数据并格式化
    for item in fixed_dungeon:
        if item["type"] == "Tower":
            result["boss_tower"].append(
                [float(item["pos"]["X"]), float(item["pos"]["Y"])]
            )
        elif item["type"] == "Fast Travel":
            result["fast_travel"].append(
                [float(item["pos"]["X"]), float(item["pos"]["Y"])]
            )

    # 将结果转换为JSON格式并保存
    with open("web/src/assets/map/points.json", "w", encoding="utf-8") as json_file:
        json.dump(result, json_file, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    # 添加命令行参数解析
    parser = argparse.ArgumentParser(description="Download tiles from palworld.gg")
    parser.add_argument(
        "--redown", action="store_true", help="Redownload existing files"
    )
    args = parser.parse_args()

    # 运行异步下载
    asyncio.run(download_images_async(args.redown))

    # # 解析JavaScript文件
    # parse_js_file()
