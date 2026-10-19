// Package mail envía correos por SMTP. Si no hay SMTP configurado se usa
// LogMailer, que solo imprime el código en la consola (útil en desarrollo).
package mail

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"log"
	"mime"
	"net"
	netmail "net/mail"
	"net/smtp"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string // ej. "Musical Game <no-reply@midominio.com>"
}

// ---------------------------------------------------------------------
// Modo desarrollo: no envía nada, imprime el código en la terminal.
// ---------------------------------------------------------------------

type LogMailer struct{}

func (LogMailer) SendPasswordResetCode(_ context.Context, to, _ string, code string, ttlMinutes int) error {
	log.Printf("[CORREO SIMULADO] Para: %s | Código de recuperación: %s (expira en %d min)", to, code, ttlMinutes)
	return nil
}

// ---------------------------------------------------------------------
// SMTP real
// ---------------------------------------------------------------------

type SMTPMailer struct {
	cfg Config
}

func NewSMTPMailer(cfg Config) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) SendPasswordResetCode(ctx context.Context, to, username, code string, ttlMinutes int) error {
	subject := "Tu código para recuperar tu contraseña"

	plain := fmt.Sprintf(`Hola %s,

Tu código para restablecer tu contraseña es:

    %s

Expira en %d minutos. Si no lo solicitaste, ignora este correo: tu contraseña no cambiará.

— Musical Game
`, username, code, ttlMinutes)

	htmlBody := fmt.Sprintf(`<div style="font-family:Arial,sans-serif;max-width:420px;margin:auto;padding:24px;color:#222">
  <h2 style="margin-top:0">Recupera tu contraseña</h2>
  <p>Hola %s, usa este código para restablecer tu contraseña:</p>
  <p style="font-size:34px;letter-spacing:10px;font-weight:bold;text-align:center;background:#f3f4f8;border-radius:8px;padding:16px">%s</p>
  <p>Expira en %d minutos.</p>
  <p style="color:#777;font-size:13px">Si no lo solicitaste, ignora este correo: tu contraseña no cambiará.</p>
</div>`, html.EscapeString(username), code, ttlMinutes)

	fromAddr, err := netmail.ParseAddress(m.cfg.From)
	if err != nil {
		return fmt.Errorf("MAIL_FROM inválido: %w", err)
	}
	to = strings.NewReplacer("\r", "", "\n", "").Replace(to) // defensa contra inyección de cabeceras

	msg := buildMessage(fromAddr.String(), to, subject, plain, htmlBody, m.cfg.Host)
	return m.send(ctx, fromAddr.Address, to, msg)
}

func (m *SMTPMailer) send(ctx context.Context, from, to string, msg []byte) error {
	addr := net.JoinHostPort(m.cfg.Host, m.cfg.Port)

	deadline := time.Now().Add(30 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	var conn net.Conn
	var err error
	if m.cfg.Port == "465" { // TLS implícito
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: m.cfg.Host})
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return err
	}
	_ = conn.SetDeadline(deadline)

	c, err := smtp.NewClient(conn, m.cfg.Host)
	if err != nil {
		conn.Close()
		return err
	}
	defer c.Close()

	// STARTTLS (puerto 587): se negocia si el servidor lo ofrece
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: m.cfg.Host}); err != nil {
			return err
		}
	}

	if m.cfg.User != "" {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("el servidor SMTP no ofrece autenticación pero se configuró SMTP_USER")
		}
		if err := c.Auth(smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)); err != nil {
			return err
		}
	}

	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

// buildMessage arma un correo multipart/alternative (texto + HTML) en UTF-8,
// con cuerpos en base64 para que las tildes lleguen bien.
func buildMessage(from, to, subject, plain, htmlBody, host string) []byte {
	idBytes := make([]byte, 12)
	_, _ = rand.Read(idBytes)
	boundary := "mg-" + hex.EncodeToString(idBytes)

	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + mime.QEncoding.Encode("UTF-8", subject) + "\r\n")
	b.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	b.WriteString("Message-ID: <" + hex.EncodeToString(idBytes) + "@" + host + ">\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n\r\n")

	writePart := func(contentType, body string) {
		b.WriteString("--" + boundary + "\r\n")
		b.WriteString("Content-Type: " + contentType + "; charset=UTF-8\r\n")
		b.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		b.WriteString(wrap76(base64.StdEncoding.EncodeToString([]byte(body))))
		b.WriteString("\r\n")
	}
	writePart("text/plain", plain)
	writePart("text/html", htmlBody)
	b.WriteString("--" + boundary + "--\r\n")
	return []byte(b.String())
}

func wrap76(s string) string {
	var out strings.Builder
	for len(s) > 76 {
		out.WriteString(s[:76] + "\r\n")
		s = s[76:]
	}
	out.WriteString(s)
	return out.String()
}
