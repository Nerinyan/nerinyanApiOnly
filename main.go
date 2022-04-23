package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"net/http"
	"runtime"
)

// TODO DOING DB 테이블 없으면 자동으로 생성하게
// TODO DOING 로그 디비에 넣을때 어떤 데이터 넣을지.
// TODO DOING 헤더로 프론트인지 api 인지 구분할수있게
//  	END   에러 핸들러.
//  	END   검색엔진 버그 체크하고 쿼리문 수정
//  	END   비트맵 반쵸에서 다운로드중에 클라이언트가 취소해도 서버는 계속 다운로드.
// TODO DOING 서버간 비트맵파일 해시값 비교해서 서로 다른경우 둘다 서버에서 삭제.
// TODO DOING 서버끼리 서로 비트맵파일 동기화 시킬수 있게
// TODO DOING 반쵸 비트맵 다운로드 제한 10분간 약 200건 10분 정지. (429 too many request) => 10분 내 100건 봇 감지 알고리즘
// TODO DOING 서버 자체적으로 10분당 150건 이내만 다운로드 가능하게 셋팅
// 		END	  검색 쿼리시 서버에 캐싱되어있는 비트맵인지 여부
// TODO DOING /status에 들어갈 상태값 추가.
// TODO END 반쵸에서 가져온 데이터 검색캐싱에 추가
// TODO DOING 디스코드 웹훅
func init() {
	config.LoadConfig()

	db.ConnectMaria()
}

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Pre(
		middleware.RemoveTrailingSlash(),
	)

	e.Use(
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		//middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
		middleware.Recover(),
	)

	// 서버상태 체크용 ====================================================================================================

	e.GET("/health", Route.Health)
	e.GET("/robots.txt", Route.Robots)
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"CpuThreadCount":        runtime.NumCPU(),
			"RunningGoroutineCount": runtime.NumGoroutine(),
			"apiCount":              *banchoCroller.ApiCount,
		})
	})

	// 비트맵 리스트 검색용 ================================================================================================
	e.GET("/search", Route.Search)
	e.GET("/search/beatmap/:mi", Route.SearchByBeatmapId)
	e.GET("/search/beatmapset/:si", Route.SearchByBeatmapSetId)

	// 개발중 || 테스트중 ===================================================================================================

	pterm.Info.Println("ECHO STARTED AT", config.Config.Port)
	e.Logger.Fatal(e.Start(":" + config.Config.Port))

}
