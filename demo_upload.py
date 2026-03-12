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
    if len(sys.argv) < 2:
        print("Sử dụng: python demo_upload.py <đường_dẫn_tới_ảnh>")
        print("Ví dụ: python demo_upload.py 5_3d_visualization.png")
        sys.exit(1)

    image_path = sys.argv[1]
    prompt_text = "Mô tả chi tiết bức ảnh này."

    print(f"Đang chuẩn bị gửi ảnh: {image_path}")
    base64_image = encode_image(image_path)

    # Khởi tạo Payload gửi đến Go Server (chuẩn Gemini/Vertex AI)
    payload = {
        "contents": [
            {
                "parts": [
                    {"text": prompt_text},
                    {
                        "inlineData": {
                            "mimeType": "image/png", # Dù là JPG hay PNG, hệ thống Go đang tự xử lý.
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

    print(f"Đang gửi yêu cầu tới {API_URL}...")
    try:
        response = requests.post(API_URL, headers=headers, data=json.dumps(payload))
        response.raise_for_status() # Báo lỗi nếu server trả về mã lỗi (500, 400...)
        
        result = response.json()
        
        print("\n--- Gemini Trả Lời ---")
        # Trích xuất nội dung văn bản từ kết quả trả về
        try:
            answer = result['candidates'][0]['content']['parts'][0]['text']
            print(answer)
        except (KeyError, IndexError) as e:
            print("Cấu trúc phản hồi không khớp dự kiến. Dữ liệu gốc:")
            print(json.dumps(result, indent=2))
        print("------------------------\n")

    except requests.exceptions.RequestException as e:
        print(f"Lỗi gọi API: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Chi tiết response: {e.response.text}")

if __name__ == "__main__":
    main()
