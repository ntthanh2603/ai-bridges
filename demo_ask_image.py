import io
import json
import base64
import requests
import sys
from PIL import Image as PILImage

# Configuration
API_URL = "http://localhost:4981/gemini/v1beta/models/gemini-advanced:generateContent"

def encode_image(image_path, max_size=(1024, 1024), quality=80):
    """Đọc file ảnh, nén/resize xuống và mã hóa thành Base64"""
    try:
        # Mở ảnh bằng Pillow
        img = PILImage.open(image_path)
        
        # Chuyển đổi sang RGB nếu là RGBA (tránh lỗi khi lưu JPEG)
        if img.mode in ("RGBA", "P"):
            img = img.convert("RGB")
            
        # Resize nếu ảnh quá lớn (giữ tỉ lệ)
        img.thumbnail(max_size, PILImage.Resampling.LANCZOS)
        
        # Lưu vào bộ nhớ đệm dạng byte với định dạng JPEG để nén dung lượng cao
        buffer = io.BytesIO()
        img.save(buffer, format="JPEG", quality=quality, optimize=True)
        
        return base64.b64encode(buffer.getvalue()).decode('utf-8')
    except ImportError:
        print("Lỗi: Bạn cần cài đặt thư viện Pillow để nén ảnh. Chạy lệnh: pip install Pillow")
        sys.exit(1)
    except Exception as e:
        print(f"Lỗi khi xử lý ảnh: {e}")
        sys.exit(1)

def main():
    # Cần ít nhất 2 tham số: tên script, đường dẫn ảnh, và câu hỏi
    if len(sys.argv) < 3:
        print("Sử dụng: python3 demo_ask_image.py <đường_dẫn_tới_ảnh> \"<câu_hỏi_của_bạn>\"")
        print("Ví dụ:   python3 demo_ask_image.py 5_3d_visualization.png \"Trục X đại diện cho cái gì?\"")
        sys.exit(1)

    image_path = sys.argv[1]
    prompt_text = sys.argv[2] # Câu hỏi từ người dùng

    print(f"Bức ảnh: {image_path}")
    print(f"Câu hỏi: {prompt_text}")
    print("Đang xủ lý và tải ảnh lên...")
    
    base64_image = encode_image(image_path)

    # Khởi tạo Payload gửi đến Go Server
    payload = {
        "contents": [
            {
                "parts": [
                    {"text": prompt_text},
                    {
                        "inlineData": {
                            "mimeType": "image/jpeg", # Định dạng ảnh chung
                            "data": base64_image
                        }
                    }
                ]
            }
        ]
    }

    headers = {
        "Content-Type": "application/json"
    }

    print(f"Đang chờ Gemini trả lời...\n")
    try:
        response = requests.post(API_URL, headers=headers, data=json.dumps(payload))
        response.raise_for_status() 
        
        result = response.json()
        
        print("============== GEMINI TRẢ LỜI ==============")
        try:
            answer = result['candidates'][0]['content']['parts'][0]['text']
            print(answer)
        except (KeyError, IndexError) as e:
            print("Cấu trúc phản hồi không khớp dự kiến. Dữ liệu gốc:")
            print(json.dumps(result, indent=2))
        print("===========================================\n")

    except requests.exceptions.RequestException as e:
        print(f"Lỗi gọi API: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Chi tiết: {e.response.text}")

if __name__ == "__main__":
    main()
