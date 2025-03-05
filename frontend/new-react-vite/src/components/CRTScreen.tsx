import "./CRTScreen.css";

function CRTScreen({ children }: { children: React.ReactNode }) {
    return (
        <div className="crt-screen">
          <div className="crt-content">
            {children}
          </div>
        </div>
    );
  }

export default CRTScreen;