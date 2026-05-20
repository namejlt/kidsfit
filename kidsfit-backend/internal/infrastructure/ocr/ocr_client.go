package ocr

import (
	"context"
	"fmt"
)

// EyeResult 单眼验光结果，包含球镜、柱镜、轴位和矫正视力
type EyeResult struct {
	// SPH 球镜度数（近视为负，远视为正）
	SPH float64 `json:"sph"`
	// CYL 柱镜度数（散光度数）
	CYL float64 `json:"cyl"`
	// AXIS 散光轴位（0-180度）
	AXIS float64 `json:"axis"`
	// VA 矫正视力（小数视力表）
	VA float64 `json:"va"`
}

// PrescriptionResult 验光单识别结果，包含双眼数据和识别置信度
type PrescriptionResult struct {
	// RightEye 右眼验光数据
	RightEye EyeResult `json:"right_eye"`
	// LeftEye 左眼验光数据
	LeftEye EyeResult `json:"left_eye"`
	// Confidence OCR识别置信度（0-1）
	Confidence float64 `json:"confidence"`
}

// OCRClient OCR识别客户端，调用OCR服务识别验光单
type OCRClient struct {
	// ocrServiceURL OCR服务的URL地址
	ocrServiceURL string
}

// NewOCRClient 创建OCR识别客户端实例
// ocrServiceURL: OCR服务的URL地址
func NewOCRClient(ocrServiceURL string) *OCRClient {
	return &OCRClient{ocrServiceURL: ocrServiceURL}
}

// RecognizePrescription 识别验光单图像，提取屈光数据
// ctx: 上下文，imageData: 验光单图像数据（JPEG/PNG格式）
// 返回验光单识别结果，包含双眼数据和置信度
// TODO: 对接真实OCR服务，当前返回模拟数据
func (c *OCRClient) RecognizePrescription(ctx context.Context, imageData []byte) (*PrescriptionResult, error) {
	// TODO: 实现真实的OCR服务调用
	// 当前返回模拟数据用于开发和测试
	if len(imageData) == 0 {
		return nil, fmt.Errorf("图像数据不能为空")
	}

	// 模拟验光单识别结果
	return &PrescriptionResult{
		RightEye: EyeResult{
			SPH:  -1.25,
			CYL:  -0.50,
			AXIS: 90,
			VA:   1.0,
		},
		LeftEye: EyeResult{
			SPH:  -1.50,
			CYL:  -0.75,
			AXIS: 85,
			VA:   0.8,
		},
		Confidence: 0.92,
	}, nil
}
