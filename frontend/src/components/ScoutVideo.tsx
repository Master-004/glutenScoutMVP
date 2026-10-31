import { useEffect, useRef } from "react";
import { trackEvent } from "../lib/clarity";

// A short Scout clip over the app. It's on a white card because the videos have
// white backgrounds, which would show as a box in dark mode. Always muted and
// skippable; if a clip can't play or motion is turned down, it steps aside.
export function ScoutVideo({ src, caption, onDone }: { src: string; caption: string; onDone(): void }) {
  const video = useRef<HTMLVideoElement>(null);
  const skip = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (window.matchMedia?.("(prefers-reduced-motion: reduce)").matches) {
      onDone();
      return;
    }
    skip.current?.focus();
    video.current?.play().catch(onDone);
  }, [onDone]);

  return (
    <div className="video-overlay" role="dialog" aria-modal="true" aria-label={caption}>
      <div className="video-card">
        <video ref={video} src={src} muted playsInline autoPlay onEnded={onDone} onError={onDone} aria-hidden="true" />
        <p className="video-caption">{caption}</p>
        <button ref={skip} type="button" className="video-skip" onClick={() => { trackEvent("video_skipped"); onDone(); }}>Skip</button>
      </div>
    </div>
  );
}
