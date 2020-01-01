package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"

	//"log"
	"math/big"
	"net"
	"time"
)

type TlsConfig struct {
	CA   string
	Cert string
	Key  string
}

func TlsConfigClient(v *TlsConfig, servername string) (*tls.Config, error) {
	if v != nil {
		return clientTlsConfig(v.CA, v.Cert, v.Key, servername)
	}

	certs, err := newTempCertificates()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ServerName:         servername,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{certs},
	}, nil
}

func TlsConfigServer(v *TlsConfig) (*tls.Config, error) {
	if v != nil {
		return serverTlsConfig(v.CA, v.Cert, v.Key)
	}

	certs, err := newTempCertificates()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certs},
		ClientAuth:   tls.RequestClientCert,
	}, nil
}

func clientTlsConfig(ca, cert, key string, addr string) (*tls.Config, error) {

	//服务端证书池
	var pool *x509.CertPool

	if ca != "" {
		//这里读取的是根证书
		buf, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, err
		}
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(buf)
	}

	//加载客户端证书
	//这里加载的是服务端签发的
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	var bSkipVerify bool

	// 如果没有配置服务端根证书，则忽略校验服务端证书有效性。
	if pool == nil {
		bSkipVerify = true
	}

	return &tls.Config{
		ServerName:         addr,
		InsecureSkipVerify: bSkipVerify,
		RootCAs:            pool,
		Certificates:       []tls.Certificate{crt},
	}, nil
}

func serverTlsConfig(ca, cert, key string) (*tls.Config, error) {
	var pool *x509.CertPool

	if ca != "" {
		//这里读取的是根证书
		buf, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, err
		}
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(buf)
	}

	//加载服务端证书
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	var authtype tls.ClientAuthType

	// 如果没有配置服务端根证书，则忽略校验客户端证书有效性。
	if pool != nil {
		authtype = tls.RequireAndVerifyClientCert
	} else {
		authtype = tls.RequestClientCert
	}

	return &tls.Config{
		Certificates: []tls.Certificate{crt},
		ClientAuth:   authtype,
		ClientCAs:    pool,
	}, nil
}

func newTempCertificates() (tls.Certificate, error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	subject := pkix.Name{                         //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Organization:       []string{"benchmark co."},
		OrganizationalUnit: []string{"benchmark"},
		CommonName:         "Go enchmark Programming",
	}
	template := x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 密钥扩展用途的序列
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048) //生成一对具有指定字位数的RSA密钥

	//CreateCertificate基于模板创建一个新的证书
	//第二个第三个参数相同，则证书是自签名的
	//返回的切片是DER编码的证书
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //DER 格式

	certOut := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	//log.Printf("CERT:%s\n", certOut.String())

	keyOut := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	//log.Printf("CERT:%s\n", keyOut.String())

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}
