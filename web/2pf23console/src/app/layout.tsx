import "./globals.css";

export const metadata = {
  title: "2pf23 web console",
  description: "2pf23 web console",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
