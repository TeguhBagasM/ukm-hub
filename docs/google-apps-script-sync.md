# Sinkronisasi Real-time ke Google Spreadsheet

Fitur ini mengirimkan data pendaftar ke Google Spreadsheet **segera setelah** registrasi
berhasil disimpan di platform. Data utama tetap tersimpan di database (sehingga alur
review → accept → convert member tetap berjalan), webhook hanya menyalinnya ke spreadsheet.

## Cara Kerja

```text
Applicant submit form
        ↓
Backend simpan Registration + Answers (transaksional)
        ↓
Backend POST JSON ke submit_webhook_url (async, non-blocking)
        ↓
Google Apps Script menambahkan baris ke Spreadsheet
```

Jika webhook gagal (mis. Apps Script error), registrasi **tetap berhasil** dan hanya
dicatat di log server.

## Setup

1. Buat Google Spreadsheet baru. Beri nama header di baris pertama sesuai field form,
   contoh: `submitted_at`, `nama_lengkap`, `nim`, `email`, `whatsapp`, `divisi_pilihan`, `alasan_bergabung`, `portfolio`.
2. Buka **Extensions → Apps Script**.
3. Tempel kode berikut, sesuaikan `SHEET_NAME` dan urutan kolom dengan header spreadsheet:

```javascript
const SHEET_NAME = 'Sheet1';

// Urutan kolom yang akan ditulis ke spreadsheet.
// Sesuaikan dengan nama field (answers_map keys) pada form kamu.
const COLUMNS = [
  'submitted_at',
  'nama_lengkap',
  'nim',
  'email',
  'whatsapp',
  'divisi_pilihan',
  'alasan_bergabung',
  'portfolio'
];

function doPost(e) {
  try {
    const payload = JSON.parse(e.postData.contents);
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(SHEET_NAME);
    const map = payload.answers_map || {};

    const row = COLUMNS.map(function (key) {
      if (key === 'submitted_at') return payload.submitted_at || '';
      return map[key] !== undefined ? map[key] : '';
    });

    sheet.appendRow(row);

    return ContentService
      .createTextOutput(JSON.stringify({ ok: true }))
      .setMimeType(ContentService.MimeType.JSON);
  } catch (err) {
    return ContentService
      .createTextOutput(JSON.stringify({ ok: false, error: String(err) }))
      .setMimeType(ContentService.MimeType.JSON);
  }
}
```

4. Klik **Deploy → New deployment**.
5. Pilih type **Web app**.
6. Set **Execute as: Me** dan **Who has access: Anyone**.
7. Klik **Deploy**, lalu salin **Web app URL** (berformat `https://script.google.com/macros/s/.../exec`).
8. Di platform, buka form event tersebut melalui API/UI admin dan isi `submit_webhook_url`
   dengan URL tersebut (wajib `https://`).

## Payload yang Dikirim Backend

```json
{
  "event": {
    "id": "uuid",
    "name": "Open Recruitment CSA 2026",
    "slug": "open-recruitment-csa-2026",
    "organization": { "id": "uuid", "name": "CSA", "slug": "creative-student-association" }
  },
  "registration_id": "uuid",
  "status": "PENDING",
  "submitted_at": "2026-09-17T12:30:49Z",
  "answers": [
    { "field_id": "uuid", "label": "Nama Lengkap", "name": "nama_lengkap", "type": "TEXT", "value": "Teguh Bagas" }
  ],
  "answers_map": {
    "nama_lengkap": "Teguh Bagas",
    "nim": "22110397",
    "divisi_pilihan": "Programming"
  }
}
```

Gunakan `answers_map` untuk pemetaan sederhana (key = `name` field), atau `answers`
jika butuh label/tipe field.

## Catatan

- Kolom spreadsheet mengikuti field dinamis. Jika admin mengubah field form, sesuaikan
  `COLUMNS` di Apps Script.
- Untuk field `CHECKBOX`, nilai dikirim sebagai string yang dipisahkan koma.
- URL webhook divalidasi harus `https://` saat dikonfigurasi admin.
