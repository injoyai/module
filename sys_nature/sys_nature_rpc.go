package sys_nature

/*



 */

type RPC struct {
	Data []Nature
}

func (this *RPC) Bytes() []byte {
	data := []byte{0x7E}            //帧头
	data = append(data, 0x00, 0x00) //数据长度,2字节
	data = append(data, 0x00)       //状态字节,预留
	for _, v := range this.Data {
		_ = v
	}

	return data
}
