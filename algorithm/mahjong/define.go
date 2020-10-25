// 麻将算法
package mahjong

// 麻将接口
type MahJong interface {
	Type() int  //牌类型; 万，筒，条，顺 ...
	Point() int //牌数值;
}
