import "./cyberpunkWindow.css";

function CyberpunkWindow({ children }: { children: React.ReactNode }){
  return (
      <div className="cyberpunk-window">
          {children}
    </div>
  );
}

export default CyberpunkWindow;