import os
import requests
import argparse
import time
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


def download_images(redown=False):
    # 计算总图片数
    total_images = sum(
        (x_max + 1) * (y_max + 1) for x_max, y_max in z_to_range.values()
    )

    # 初始化 tqdm 进度条
    progress_bar = tqdm(total=total_images, desc="Downloading images", unit="img")

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

                # 检查本地文件是否已经存在
                if os.path.exists(file_path) and not redown:
                    progress_bar.update(1)
                    continue

                # 尝试下载图片
                attempt = 0
                success = False
                while attempt < RETRY_TIMES:
                    try:
                        response = requests.get(url, headers=headers)
                        # 检查响应状态码
                        if response.status_code == 200:
                            # 将图片内容保存到本地
                            with open(file_path, "wb") as f:
                                f.write(response.content)
                            success = True
                            break  # 成功下载后退出重试循环
                        elif response.status_code == 404:
                            print(f"Skipped {url} - Not Found (404)")
                            break  # 404 不重试
                        elif response.status_code == 403:
                            print(f"Skipped {url} - Forbidden (403)")
                            break  # 403 不重试
                        else:
                            print(
                                f"Failed to download {url} (status code: {response.status_code})"
                            )
                            break  # 其他错误码不重试
                    except Exception as e:
                        print(f"Error downloading {url}: {e}")
                    attempt += 1
                    if attempt < RETRY_TIMES:
                        print(f"Retrying ({attempt}/{RETRY_TIMES}) for {url}")
                        time.sleep(1)  # 重试前等待 1 秒

                # 更新进度条
                if success:
                    progress_bar.update(1)
                else:
                    print(f"Failed to download {url} after {RETRY_TIMES} attempts")
                    progress_bar.update(1)  # 即使失败也更新进度条

    progress_bar.close()


if __name__ == "__main__":
    # 添加命令行参数解析
    parser = argparse.ArgumentParser(description="Download tiles from palworld.gg")
    parser.add_argument(
        "--redown", action="store_true", help="Redownload existing images"
    )
    args = parser.parse_args()

    # 根据参数调用下载函数
    download_images(redown=args.redown)
