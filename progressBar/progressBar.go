package progressBar

import (
	"fmt"
)

type ProgressBar struct {
	percent int //百分比
	cur  int  //当前进度
	total int //总进度
	rate  string //进度条
	graph string //显示字符
}

func New(start, total int) *ProgressBar {
	var pb ProgressBar
	pb.total = total
	pb.cur   = start
	pb.percent = pb.GetPercentage()
	pb.graph = "█"
	for i := 0;i< pb.percent/2; i++ {
		pb.rate +=pb.graph
	}
	return &pb
}

func (p *ProgressBar) Draw(cur int) {
	p.cur = cur
	last := p.percent
	p.percent = p.GetPercentage()
	if p.percent != last {
		i  := p.percent-last
		if p.percent %2 == 1 && last %2 == 0 {
			i = i-1
		}else if p.percent %2 == 0 && last %2 == 1 {
			i = i+1
		}
		for j :=0 ;j< i/2;j++ {
			p.rate += p.graph
		}
	}
	//\r表示回到行初始 - 是左对齐的意思，默认是右对齐 %% 输出%号
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", p.rate, p.percent, p.cur, p.total)
}

func (p *ProgressBar) GetPercentage() int{
	return int(float64(p.cur)/float64(p.total)*100)
}

func (p *ProgressBar) Finish(){
	fmt.Println()
}

