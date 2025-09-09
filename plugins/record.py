import json
import time
from pynput import mouse, keyboard
from threading import Lock

class MouseRecorder:
    def __init__(self):
        self.current_recording = []
        self.is_recording = False
        self.lock = Lock()
        self.raw_config_file = "mouse_movements_raw.json"
        self.merged_config_file = "mouse_movements_merged.json"

        self.mouse_listener = mouse.Listener(on_move=self.on_move, on_click=self.on_click)
        self.keyboard_listener = keyboard.Listener(on_release=self.on_key_release)

        self.mouse_listener.start()
        self.keyboard_listener.start()

    def on_move(self, x, y):
        if self.is_recording and self.current_recording:
            with self.lock:
                current_time = time.time()
                last_time = self.current_recording[-1]["time"] if self.current_recording else current_time
                self.current_recording.append({
                    "dx": x - self.current_recording[-1]["x"] if self.current_recording else 0,
                    "dy": y - self.current_recording[-1]["y"] if self.current_recording else 0,
                    "x": x,
                    "y": y,
                    "time": current_time,
                    "relative_time": current_time - last_time
                })

    def on_click(self, x, y, button, pressed):
        if button == mouse.Button.left:
            if pressed:
                with self.lock:
                    self.is_recording = True
                    self.current_recording = [{
                        "dx": 0,
                        "dy": 0,
                        "x": x,
                        "y": y,
                        "time": time.time(),
                        "relative_time": 0
                    }]
                    print("Recording started")
            else:
                with self.lock:
                    self.is_recording = False
                    if self.current_recording:
                        self.save_to_config()
                        print("Recording saved to both raw and merged files.")

    def on_key_release(self, key):
        try:
            if key == keyboard.Key.enter:
                with self.lock:
                    print("Replaying last recording...")
                    self.replay(self.current_recording)
            elif key == keyboard.Key.esc:
                print("Exiting program...")
                self.mouse_listener.stop()
                self.keyboard_listener.stop()
                return False
        except AttributeError:
            pass

    def merge_consecutive_movements(self, data):
        """合并连续相同方向的鼠标移动事件"""
        if not data:
            return []

        merged_data = []
        current_dx = data[0]['dx']
        current_dy = data[0]['dy']
        total_dx = data[0]['dx']
        total_dy = data[0]['dy']
        total_time = data[0]['relative_time']

        for i in range(1, len(data)):
            if data[i]['dx'] == current_dx and data[i]['dy'] == current_dy:
                total_dx += data[i]['dx']
                total_dy += data[i]['dy']
                total_time += data[i]['relative_time']
            else:
                merged_data.append({
                    'dx': total_dx,
                    'dy': total_dy,
                    'relative_time': total_time
                })
                current_dx = data[i]['dx']
                current_dy = data[i]['dy']
                total_dx = data[i]['dx']
                total_dy = data[i]['dy']
                total_time = data[i]['relative_time']

        merged_data.append({
            'dx': total_dx,
            'dy': total_dy,
            'relative_time': total_time
        })

        return merged_data

    def save_to_config(self):
        """保存记录到配置文件，同时输出原始和合并后的数据"""
        # 首先过滤出需要的字段
        filtered_recording = [
            {
                "dx": point["dx"],
                "dy": point["dy"],
                "relative_time": point["relative_time"]
            }
            for point in self.current_recording
        ]

        # 保存原始数据
        with open(self.raw_config_file, "w") as f:
            json.dump(filtered_recording, f, indent=4)

        # 合并连续相同方向的移动
        merged_recording = self.merge_consecutive_movements(filtered_recording)

        # 保存合并后的数据
        with open(self.merged_config_file, "w") as f:
            json.dump(merged_recording, f, indent=4)

    def load_from_config(self):
        try:
            # 默认加载合并后的配置
            with open(self.merged_config_file, "r") as f:
                self.current_recording = json.load(f)
            print("Loaded merged recording")
        except FileNotFoundError:
            try:
                # 如果合并文件不存在，尝试加载原始文件
                with open(self.raw_config_file, "r") as f:
                    self.current_recording = json.load(f)
                print("Loaded raw recording")
            except (FileNotFoundError, json.JSONDecodeError):
                self.current_recording = []
                print("No recording found")

    def replay(self, recording):
        if not recording:
            print("No recording to replay.")
            return

        mouse_controller = mouse.Controller()

        # 初始移动
        mouse_controller.move(recording[0]["dx"], recording[0]["dy"])

        for point in recording[1:]:
            time.sleep(point["relative_time"])
            mouse_controller.move(point["dx"], point["dy"])

if __name__ == "__main__":
    recorder = MouseRecorder()
    recorder.load_from_config()

    recorder.mouse_listener.join()
    recorder.keyboard_listener.join()