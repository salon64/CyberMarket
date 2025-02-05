import "./cyberpunkWindow.css";

function cyberpunkWindow({ children }: { children: React.ReactNode }){
  return (
      <div className="cyberpunk-window">
          {children}
    </div>
  );
}

export default cyberpunkWindow;