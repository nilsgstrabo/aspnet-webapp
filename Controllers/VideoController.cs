using System.IO;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Microsoft.AspNetCore.StaticFiles;
using Microsoft.Graph;
using System.Net.Http;
using System.Net.Http.Headers;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Http;

namespace aspnet_webapp.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class VideoController : ControllerBase
    {
        private readonly ILogger _logger;
        private readonly IConfiguration _configuration;

        public VideoController(IConfiguration configuration,ILogger<VideoController> logger)
        {
            _logger = logger;
            _configuration = configuration;
        }

        [HttpGet()]
        public IActionResult GetVideos()
        {
            
            return this.Ok();
        }

        [HttpGet("{name}")]
        public IActionResult GetVideo(string name)
        {
            var videoPath=_configuration["VIDEO_PATH"];
            if (videoPath?.Length==0) {
                return this.StatusCode(500);
            }
            _logger.LogInformation("Stream video {0}, range {1}", name, this.Request.GetTypedHeaders().Range?.ToString());
            var filename=System.IO.Path.Combine(_configuration["VIDEO_PATH"], new System.IO.FileInfo(name).Name);
            return this.PhysicalFile(filename, "video/mp4", true);
        }

    }
}