import { ImageResponse } from "next/og";

export const runtime = "edge";

export const alt = "Lock In";
export const size = {
  width: 1200,
  height: 630,
};

export const contentType = "image/png";

export default function og() {
  return new ImageResponse(
    (
      <div
        style={{
          backgroundColor: "#0f172a",
          backgroundSize: "150px 150px",
          height: "100%",
          width: "100%",
          display: "flex",
          textAlign: "center",
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          flexWrap: "nowrap",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            justifyItems: "center",
          }}
        >
          <span
            style={{
              fontSize: 130,
              fontStyle: "normal",
              letterSpacing: "-0.025em",
              color: "white",
              marginTop: 30,
              padding: "0 120px",
              lineHeight: 1.4,
              whiteSpace: "pre-wrap",
            }}
          >
            Lock In
          </span>
        </div>
      </div>
    ),
    {
      ...size,
    },
  );
}
