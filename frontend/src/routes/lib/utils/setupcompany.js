// utils/setupcompany.js

import {
    userDataSession,
    setOrganizationSession,
  } from "$lib/utils/auth";

export async function emailBlast(email, password) {
    try {
      await axios.post('http://localhost:3000/api/send-email', {
        subject: 'Pemberitahuan Pendaftaran Akun di BJJ Management Asset',
        message: `
  Selamat!
  
  Akun Anda telah berhasil didaftarkan di BJJ Management Asset, platform terdepan yang menyederhanakan proses penandatanganan dokumen secara elektronik.
  
  Berikut adalah informasi akun Anda:
  Email: ${email}
  Password: ${password}
  
  Bergabung dengan BJJ Management Asset akan memberdayakan Anda untuk dengan mudah menandatangani dan mengelola dokumen secara digital, meningkatkan efisiensi proses Anda.
  
  Jika Anda memiliki pertanyaan atau memerlukan bantuan selama proses pendaftaran, jangan ragu untuk menghubungi tim dukungan kami di [Alamat Email Dukungan].
  
  Terima kasih telah memilih BJJ Management Asset untuk kebutuhan tanda tangan digital Anda. Kami berharap dapat menyambut Anda di komunitas kami!
  
  Salam terbaik,
  [Tanda Tangan Anda]
  [Tautan Pusat Bantuan atau Kontak Dukungan]
  `,
        to: email,
      });
  
    } catch (error) {
      console.error('Error sending email:', error);
    }
  }  

  export async function checkOrganization() {
    const result = userDataSession();
    const response2 = await fetch(
        `/api/auth?jwt=${result.jwt}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      if (response2.ok) {
        const result = await response2.json();
        if (result.organization !== null){
          setOrganizationSession(result.organization);
          return true;
        } else {
          return false;
        }
      }
  }