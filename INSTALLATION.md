# Go Kurulum ve Test Rehberi

## 📥 Adım 1: Go'yu İndir ve Kur

1. **İndirme Sayfası Açıldı**: https://go.dev/dl/
2. **Windows için en son sürümü indir**: `go1.21.x.windows-amd64.msi` (veya daha yeni)
3. **İndirilen .msi dosyasını çalıştır**
4. **Kurulum sihirbazını takip et** (varsayılan ayarlar yeterli)
5. **Kurulum tamamlandıktan sonra terminal'i KAPAT ve YENİDEN AÇ**

## ✅ Adım 2: Go Kurulumunu Doğrula

Yeni bir PowerShell penceresi aç ve şunu çalıştır:

```powershell
go version
```

Çıktı şöyle olmalı:
```
go version go1.21.x windows/amd64
```

## 🧪 Adım 3: Projeyi Test Et

### Testleri Çalıştır

```powershell
cd "C:\Users\ASUS\OneDrive\Masaüstü\Backend-Systems Focus\redis-clone"

# Bağımlılıkları indir
go mod download

# Tüm testleri çalıştır
go test ./... -v

# Sadece store testleri
go test ./internal/store -v

# Coverage ile
go test ./... -cover
```

### Benchmarkları Çalıştır

```powershell
go test -bench=. -benchmem ./internal/store
```

### Build Et

```powershell
# Windows için build
go build -o redis-clone.exe ./cmd/server

# Çalıştır
.\redis-clone.exe
```

## 🔌 Adım 4: Sunucuyu Test Et

### Terminal 1: Sunucuyu Başlat
```powershell
.\redis-clone.exe
```

Şunu göreceksiniz:
```
╔═══════════════════════════════════════════════════════════╗
║                    REDIS CLONE                            ║
╚═══════════════════════════════════════════════════════════╝
🚀 Redis Clone server started on :6379
📊 Ready to accept connections...
```

### Terminal 2: Telnet ile Bağlan

Yeni bir PowerShell penceresi aç:

```powershell
# Telnet'i etkinleştir (ilk kez)
Enable-WindowsOptionalFeature -Online -FeatureName TelnetClient

# Bağlan
telnet localhost 6379
```

### Komutları Test Et

Telnet bağlandıktan sonra:

```
PING
# Yanıt: +PONG

SET name "Redis Clone"
# Yanıt: +OK

GET name
# Yanıt: $11
#        Redis Clone

SET session "xyz123" EX 60
# Yanıt: +OK

TTL session
# Yanıt: :60

KEYS *
# Yanıt: *2
#        $4
#        name
#        $7
#        session

INFO
# Yanıt: Server bilgileri
```

## 🐛 Sorun Giderme

### "go: command not found"
- Terminal'i kapat ve yeniden aç
- PATH'i kontrol et: `$env:PATH`

### "cannot find package"
- `go mod download` çalıştır
- `go mod tidy` dene

### Port zaten kullanımda
- Farklı port kullan: `.\redis-clone.exe -addr :6380`

### Testler fail oluyor
- Go version kontrol et: `go version` (1.21+ olmalı)
- Modülleri güncelle: `go mod tidy`

## 📊 Beklenen Test Sonuçları

Tüm testler geçmeli:
```
ok      github.com/yourusername/redis-clone/internal/store      0.234s  coverage: 85.2%
ok      github.com/yourusername/redis-clone/internal/protocol   0.156s  coverage: 92.1%
ok      github.com/yourusername/redis-clone/internal/commands   0.189s  coverage: 88.5%
ok      github.com/yourusername/redis-clone/internal/server     0.312s  coverage: 81.3%
```

## 🎯 Başarı Kriterleri

✅ `go version` çalışıyor
✅ `go test ./...` tüm testler geçiyor
✅ `go build` başarılı
✅ Sunucu başlıyor
✅ Telnet ile bağlanabiliyor
✅ Komutlar çalışıyor

## 💡 İpuçları

- **Ctrl+C** ile sunucuyu durdur
- **Ctrl+]** sonra `quit` ile telnet'ten çık
- Test sırasında farklı portlar kullan (16379, 16380, vb.)

---

**Kurulum tamamlandıktan sonra bu dosyayı referans olarak kullanabilirsiniz!**
